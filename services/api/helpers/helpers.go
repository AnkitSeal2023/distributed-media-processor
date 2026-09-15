package helpers

import (
	"log"
	"net/http"

	"distributed-media-processing-platform/constants/error_msgs"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
