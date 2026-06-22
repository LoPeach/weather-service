package response

import (
	"encoding/json"
	"net/http"
	"new_begin/internal/models"
)

func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "json encode error", http.StatusInternalServerError)
	}
}

func JSON(w http.ResponseWriter, status int, resp models.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(resp)
}

func Success(w http.ResponseWriter, status int, data any) {
	JSON(w, status, models.Response{
		Data: data,
	})
}

func Error(w http.ResponseWriter, status int, message string) {
	JSON(w, status, models.Response{
		Error: message,
	})
}
