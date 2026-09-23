package repository

import (
	"context"
	"distributed-media-processing-platform/constants/error_msgs"
	"errors"
	"log"

	"github.com/google/uuid"
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
func (r *UploadRepository) InsertNewVideo(ctx context.Context, filename string, userid string, status string) error {
	userID, err := uuid.Parse(userid)
	if err != nil {
		log.Printf("ERR: Failed to parse userID:%v", err)
		return error_msgs.ErrInternalServer
	}

	log.Printf("userid:%v\n", userID)
	_, err = r.db.Exec(ctx, InsertQuery, filename, userid, status)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return errors.New("A video with such filename already exists")
		} else {
			log.Printf("ERR: Failed to create user:%v", err)
			return error_msgs.ErrInternalServer
		}
	}
	return nil
}
