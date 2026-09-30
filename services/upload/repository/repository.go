package repository

import (
	"context"
	"distributed-media-processing-platform/constants/error_msgs"
	"errors"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UploadRepository struct {
	db *pgxpool.Pool
}

func NewUploadRepository(pool *pgxpool.Pool) *UploadRepository {
	return &UploadRepository{
		db: pool,
	}
}
func (r *UploadRepository) InsertNewVideo(ctx context.Context, filename string, userid string, status string) (vid_id string, e error) {

	var videoID string
	userID, err := uuid.Parse(userid)
	if err != nil {
		log.Printf("ERR: Failed to parse userID:%v", err)
		return "", error_msgs.ErrInternalServer
	}

	row := r.db.QueryRow(ctx, InsertVideo, filename, userID, status)
	err = row.Scan(&videoID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.Is(err, pgx.ErrNoRows) {
			return "", error_msgs.ErrUnauthorized
		} else if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", errors.New("23505")
		} else {
			log.Printf("ERR occurred while inserting video: %v", err)
			return "", error_msgs.ErrInternalServer
		}
	}
	return videoID, nil
}

func (r *UploadRepository) UpdateVideoStatus(ctx context.Context, videoID string, status string) error {
	split := strings.Split(videoID, ".")
	if len(split) < 2 {
		log.Printf("ERR: Invalid videoID format: %v", videoID)
		return errors.New("Bad video format")
	}
	vidID, err := uuid.Parse(split[0])
	if err != nil {
		log.Printf("ERR: Failed to parse videoID:%v video id: %v", err, videoID)
		return error_msgs.ErrInternalServer
	}

	cmdTag, err := r.db.Exec(ctx, UpdateVideoStatus, status, vidID)
	// TODO: handle db errors
	if err != nil {
		return error_msgs.ErrInternalServer
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("videoID does not exist")
	}
	return nil
}
