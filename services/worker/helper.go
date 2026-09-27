package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v7"
)

func processVideo(ctx context.Context, fileNameWithExtension string) (e error) {
	videoID := strings.Split(fileNameWithExtension, ".")[0]

	ok := downloadVideo(ctx, fileNameWithExtension)
	if !ok {
		return fmt.Errorf("failed to download video with id: %v", fileNameWithExtension)
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
		return fmt.Errorf("failed to run ffprobe: %v", err)
	}

	duration, err := strconv.ParseFloat(strings.TrimSpace(string(op)), 64)
	if err != nil {
		return fmt.Errorf("failed to parse video duration: %v", err)
	}

	fmt.Printf("Video duration: %.2f seconds\n", duration)

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
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	cmd.Stderr = io.Discard
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ffmpeg: %v", err)
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

			fmt.Printf("Progress: %.2f%%\n", progress)
		}

		if line == "progress=end" {
			fmt.Println("Progress: 100%")
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed reading ffmpeg output: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("ffmpeg failed: %v", err)
	}

	fmt.Println("Transcoding completed")

	return nil
}

func downloadVideo(ctx context.Context, videoID string) (ok bool) {
	object, err := minioClient.GetObject(ctx, bucketName, videoID, minio.GetObjectOptions{})
	if err != nil {
		// TODO: remove printf
		fmt.Printf("err while fetching object: %v \n", err)
		// TODO: handle deletion of video id entry from database
		return false
	}
	defer object.Close()

	localFile, err := os.Create(videoID)
	if err != nil {
		fmt.Printf("err while creating local file: %v \n", err)
		return false
	}
	defer localFile.Close()

	if _, err = io.Copy(localFile, object); err != nil {
		fmt.Printf("err while copying object to local file: %v \n", err)
		fmt.Println(err)
		return false
	}
	return true
}
