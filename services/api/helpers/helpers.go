package helpers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"distributed-media-processing-platform/constants/error_msgs"
	"os"

	"github.com/golang-jwt/jwt/v5"
	env "github.com/joho/godotenv"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userJWTCustomClaims struct {
	UserID string `json:"user-id"`
	jwt.RegisteredClaims
}

func HandleGrpcError(w http.ResponseWriter, err error) {
	switch status.Code(err) {
	case codes.AlreadyExists:
		WriteError(w, http.StatusConflict, error_msgs.ErrUserAlreadyExists.Error())
	case codes.Unauthenticated:
		WriteError(w, http.StatusUnauthorized, error_msgs.ErrUnauthorized.Error())
	default:
		WriteError(w, http.StatusInternalServerError, error_msgs.ErrInternalServer.Error())
	}

}

func WriteError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Printf("Failed to write message")
	}
}

func VerifyAccessToken(r *http.Request) (string, error) {
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

	token, err := jwt.ParseWithClaims(accessToken, &userJWTCustomClaims{}, func(token *jwt.Token) (any, error) {
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
