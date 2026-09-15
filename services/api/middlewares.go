package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"context"

	"github.com/golang-jwt/jwt/v5"
	env "github.com/joho/godotenv"
)

type userJWTCustomClaims struct {
	UserID string `json:"user-id"`
	jwt.RegisteredClaims
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, err := verifyAccessToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func verifyAccessToken(r *http.Request) (string, error) {
	err := env.Load()
	jwt_secret := os.Getenv("ACCESS_KEY_JWT_SECRET")
	if err != nil {
		log.Printf("error occured during loading env variables: %v", err)
		return "", err
	}
	AuthorisationHeaders := r.Header.Get("Authorization")
	if AuthorisationHeaders == "" {
		return "", fmt.Errorf("Authorization headers are missing")
	}

	parts := strings.Split(AuthorisationHeaders, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("Bearer token missing from authorization header")
	}
	accessToken := parts[1]

	token, err := jwt.ParseWithClaims(accessToken, &userJWTCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwt_secret, nil
	})

	if claims, ok := token.Claims.(userJWTCustomClaims); ok && token.Valid {
		fmt.Printf("User ID from token: %s\n", claims.UserID)
		return claims.UserID, nil
	} else {
		return "", errors.New("invalid token")
	}
}
