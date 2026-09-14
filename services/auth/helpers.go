package main

import (
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

func generateAccessToken(userID string) (string, error) {
	err := env.Load()
	if err != nil {
		log.Fatalf("error occured during loading env variables: %v", err)
	}
	envJwtSecret := os.Getenv("JWT_SECRET")
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
