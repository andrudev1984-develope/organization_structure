package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"organization_structure/internal/model/out"
)

func sendError(w http.ResponseWriter, status int, message string, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	slog.Error("error is occurred",
		"status", status,
		"msg", message,
		"err", err,
	)

	json.NewEncoder(w).Encode(out.CustomError{
		Code:    status,
		Message: message,
	})
}

func sendSuccess(w http.ResponseWriter, status int, message string, content any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	slog.Info(message,
		"status", status,
	)

	if content != nil {
		json.NewEncoder(w).Encode(content)
	}
}
