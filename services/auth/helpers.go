package main

import (
	"context"
	"distributed-media-processing-platform/constants/error_msgs"
	"log"
	"os"
	"time"

	"github.com/alexedwards/argon2id"
	jwt "github.com/golang-jwt/jwt/v5"
	env "github.com/joho/godotenv"
)

func hashPassword(userPassword string) (string, error) {
	hash, err := argon2id.CreateHash(userPassword, argon2id.DefaultParams)
	if err != nil {
		log.Printf("error occured during hashing password: %v", err)
		return "", err
	}
	return hash, nil
}

func checkPasswordHash(ctx context.Context, email, userPassword string) (userid string, match bool, errors error) {
	userid, hash, err := AuthRepo.GetPasswordHashByEmail(ctx, email)
	if err != nil {
		return "", false, err
	}
	match, err = argon2id.ComparePasswordAndHash(userPassword, hash)
	if match {
		return userid, true, nil
	} else {
		if err != nil {
			log.Printf("error occured during comparing password and hash: %v", err)
			return "", false, error_msgs.ErrInternalServer
		}
		return "", false, error_msgs.ErrUnauthorized
	}
}

func generateAccessToken(userID string) (string, error) {
	err := env.Load()
	if err != nil {
		log.Fatalf("error occured during loading env variables: %v", err)
	}
	envJwtSecret := os.Getenv("ACCESS_KEY_JWT_SECRET")
	jwtSecret := []byte(envJwtSecret)
	claims := jwt.MapClaims{
		"user-id": userID,
		"exp":     time.Now().Add(14 * time.Minute).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtToken, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("error occured while generating jwt token for user %v : ERR:%v", userID, err)
		return "", err
	}

	return jwtToken, nil
}
