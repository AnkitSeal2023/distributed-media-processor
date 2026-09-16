package main

import (
	"distributed-media-processing-platform/services/api/helpers"
	"net/http"

	"context"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid, err := helpers.VerifyAccessToken(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "userID", uid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
