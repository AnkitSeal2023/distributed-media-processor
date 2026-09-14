package main

import (
	"context"
	"distributed-media-processing-platform/constants/error_msgs"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: pool,
	}
}

func (r *AuthRepository) CreateUser(ctx context.Context, email string, passwordHash string) (userID string, err error) {
	var userid uuid.UUID
	query := "INSERT INTO users(email, password_hash) VALUES ($1, $2)	RETURNING userid;"

	row := r.db.QueryRow(ctx, query, email, passwordHash)
	err = row.Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", error_msgs.ErrUserAlreadyExists
		} else {
			log.Printf("ERR: Failed to create user:%v", err)
			return "", error_msgs.ErrInternalServer
		}
	}

	return userid.String(), nil
}
