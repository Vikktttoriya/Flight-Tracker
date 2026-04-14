package dto

import (
	"encoding/json"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if body != nil {
		_ = json.NewEncoder(w).Encode(body)
	}
}

func RespondError(w http.ResponseWriter, status int, code, message string) {
	RespondJSON(w, status, map[string]string{
		"code":    code,
		"message": message,
	})
}
