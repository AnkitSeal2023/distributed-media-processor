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
		st, ok := status.FromError(err)
		if !ok {
			EncodeResponse(w, nil, error_msgs.ErrUserAlreadyExists.Error(), nil, http.StatusConflict)
			return
		}
		EncodeResponse(w, nil, st.Message(), nil, http.StatusConflict)

	case codes.Unauthenticated:
		st, ok := status.FromError(err)
		if !ok {
			EncodeResponse(w, nil, error_msgs.ErrUnauthorized.Error(), nil, http.StatusUnauthorized)
			return
		}
		EncodeResponse(w, nil, st.Message(), nil, http.StatusUnauthorized)

	case codes.InvalidArgument:
		st, ok := status.FromError(err)
		if !ok {
			EncodeResponse(w, nil, error_msgs.ErrInvalidArgument.Error(), nil, http.StatusInternalServerError)
			return
		}
		EncodeResponse(w, nil, st.Message(), nil, http.StatusBadRequest)
	default:
		fmt.Printf("err: %v\n", err.Error())
		EncodeResponse(w, nil, errors.New("an error occurred").Error(), nil, int(status.Code(err)))
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

	token, err := jwt.ParseWithClaims(accessToken, &userJWTCustomClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(accessKeyJwtSecret), nil
	})

	if err != nil {
		fmt.Printf("Error parsing token: %v\n", err)
		return "", err
	}

	if claims, ok := token.Claims.(*userJWTCustomClaims); ok && token.Valid {
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
