package main

import (
	"bufio"
	"context"
	"distributed-media-processing-platform/services/worker/repository"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v7"
)

func processVideo(ctx context.Context, fileNameWithExtension string, workerRepo repository.WorkerRepository) (e error, requeue bool) {
	videoID := strings.Split(fileNameWithExtension, ".")[0]

	ok := downloadVideo(ctx, fileNameWithExtension)
	if !ok {
		return fmt.Errorf("failed to download video with id: %v", fileNameWithExtension), true
	}

	// Get video duration
	ffprobe := exec.Command(
		"ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		fileNameWithExtension,
	)

	op, err := ffprobe.Output()
	if err != nil {
		db_err := workerRepo.UpdateVideoStatusToFailed(ctx, videoID)
		if db_err != nil {
			fmt.Printf("failed to update video status to failed: video_ID = %v \n err=%v", videoID, db_err)
			fmt.Printf("ffprobe failed for video id: %v err = %v", videoID, err)
			return fmt.Errorf("failed to update video status to failed: %v", db_err), false
		}
		return fmt.Errorf("failed to execute ffprobe: %v", err), false
	}

	duration, err := strconv.ParseFloat(strings.TrimSpace(string(op)), 64)
	if err != nil {
		return fmt.Errorf("failed to parse video duration: %v", err), false
	}

	cmd := exec.Command(
		"ffmpeg",
		"-i", fileNameWithExtension,

		"-filter_complex",
		"[0:v]split=3[v720][v480][v144];"+
			"[v720]scale=-2:720[v720out];"+
			"[v480]scale=-2:480[v480out];"+
			"[v144]scale=-2:144[v144out]",

		"-map", "[v720out]",
		"-map", "0:a?",
		"-c:v", "libx264",
		"-b:v", "2500k",
		"-c:a", "aac",
		"-b:a", "128k",
		videoID+"_720p.mp4",

		"-map", "[v480out]",
		"-map", "0:a?",
		"-c:v", "libx264",
		"-b:v", "1200k",
		"-c:a", "aac",
		"-b:a", "128k",
		videoID+"_480p.mp4",

		"-map", "[v144out]",
		"-map", "0:a?",
		"-c:v", "libx264",
		"-b:v", "150k",
		"-c:a", "aac",
		"-b:a", "64k",
		videoID+"_144p.mp4",

		"-progress", "pipe:1",
		"-nostats",
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err), true
	}

	cmd.Stderr = io.Discard
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %v", err), true
	}

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		line := scanner.Text()

		if after, ok0 := strings.CutPrefix(line, "out_time_us="); ok0 {
			value := after

			outTimeUS, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				continue
			}

			progress := float64(outTimeUS) / (duration * 1_000_000) * 100

			if progress > 100 {
				progress = 100
			}

			// TODO: remove printf
			fmt.Printf("Progress: %.2f%%\n", progress)
		}

		if line == "progress=end" {
			// TODO: remove printf
			fmt.Println("Progress: 100%")
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed reading ffmpeg output: %v", err), true
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg failed: %v", err), true
	}
	fmt.Println("Transcoding completed")

	err = workerRepo.UpdateVideoStatusToCompleted(ctx, videoID)
	if err != nil {
		return fmt.Errorf("failed to update video status to completed: %v", err), false
	}
	return nil, false
}

func downloadVideo(ctx context.Context, filename string) (ok bool) {
	err := minioClient.FGetObject(ctx, bucketName, filename, filename, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("error while downloading video:", err)
		return
	}
	return true
}
