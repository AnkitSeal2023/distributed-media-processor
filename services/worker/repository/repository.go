package repository

import (
	"context"
	"distributed-media-processing-platform/constants/error_msgs"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkerRepository struct {
	db *pgxpool.Pool
}

func NewWorkerRepository(pool *pgxpool.Pool) *WorkerRepository {
	return &WorkerRepository{
		db: pool,
	}
}

func (r *WorkerRepository) UpdateVideoStatusToCompleted(ctx context.Context, videoID string) (e error) {

	query := updateVideoStatusToCompleted
	cmd, err := r.db.Exec(ctx, query, videoID)
	if err != nil {
		return error_msgs.ErrInternalServer
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("videoID does not exist")
	}
	return nil
}

func (r *WorkerRepository) UpdateVideoStatusToFailed(ctx context.Context, videoID string) (e error) {

	query := updateVideoStatusToFailed
	cmd, err := r.db.Exec(ctx, query, videoID)
	if err != nil {
		return error_msgs.ErrInternalServer
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("videoID does not exist")
	}
	return nil
}

func (r *WorkerRepository) UpdateRetryCount(ctx context.Context, videoID string, lastError string) (e error) {
	query := updateRequeuesWithFailure
	cmd, err := r.db.Exec(ctx, query, videoID, lastError)
	if err != nil {
		return error_msgs.ErrInternalServer
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("videoID does not exist")
	}
	return nil
}
