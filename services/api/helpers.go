package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"distributed-media-processing-platform/constants/error_msgs"

	"github.com/golang-jwt/jwt/v5"

	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userJWTCustomClaims struct {
	UserID string `json:"user-id"`
	jwt.RegisteredClaims
}

type apiResponse struct {
	Data  any    `json:"data,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Extra any    `json:"extra,omitempty"`
}

func HandleGrpcError(w http.ResponseWriter, err error) {
	switch status.Code(err) {
	case codes.AlreadyExists:
		EncodeResponse(w, nil, error_msgs.ErrUserAlreadyExists.Error(), nil, http.StatusConflict)
	case codes.Unauthenticated:
		EncodeResponse(w, nil, error_msgs.ErrUnauthorized.Error(), nil, http.StatusUnauthorized)
	default:
		EncodeResponse(w, nil, error_msgs.ErrInternalServer.Error(), nil, http.StatusInternalServerError)
	}

}

func VerifyAccessToken(r *http.Request) (string, error) {
	AuthorisationHeaders := r.Header.Get("Authorization")
	if AuthorisationHeaders == "" {
		return "", fmt.Errorf("Authorization headers are missing")
	}

	parts := strings.Split(AuthorisationHeaders, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("Bearer token missing from authorization header")
	}
	accessToken := parts[1]

	token, _ := jwt.ParseWithClaims(accessToken, &userJWTCustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return accessKeyJwtSecret, nil
	})

	if claims, ok := token.Claims.(userJWTCustomClaims); ok && token.Valid {
		fmt.Printf("User ID from token: %s\n", claims.UserID)
		return claims.UserID, nil
	} else {
		return "", errors.New("invalid token")
	}
}

func EncodeResponse(w http.ResponseWriter, data any, msg string, extra any, code int) {
	w.Header().Set("Content-Type", "application/json")

	response := apiResponse{
		Data:  data,
		Msg:   msg,
		Extra: extra,
	}

	json, err := json.Marshal(response)
	if err != nil {
		log.Printf("could not encode JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)

	_, err = w.Write(json)
	if err != nil {
		log.Printf("Failed to write to response:%v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
