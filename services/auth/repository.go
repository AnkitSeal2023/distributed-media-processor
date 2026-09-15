package main

import (
	"context"
	"distributed-media-processing-platform/constants/error_msgs"
	"errors"
	"log"

	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v5"
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

func (r *AuthRepository) CreateUser(ctx context.Context, email string, passwordHash string) (userID string, error error) {
	var userid uuid.UUID
	query := "INSERT INTO users(email, password_hash) VALUES ($1, $2)	RETURNING userid;"

	row := r.db.QueryRow(ctx, query, email, passwordHash)
	err := row.Scan(&userid)
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

func (r *AuthRepository) GetPasswordHashByEmail(ctx context.Context, email string) (userID string, passwordHash string, error error) {

	var user_hash_password string
	var userid uuid.UUID

	query := "SELECT u.password_hash, u.userid  FROM users u WHERE u.email = $1;"
	row := r.db.QueryRow(ctx, query, email)
	err := row.Scan(&user_hash_password, &userid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", error_msgs.ErrUnauthorized
		} else {
			log.Printf("ERR occurred while searching for user's password hash: %v", err)
			return "", "", error_msgs.ErrInternalServer
		}
	}

	return userid.String(), user_hash_password, nil

}
