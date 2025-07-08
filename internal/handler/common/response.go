package common

import (
	"encoding/json"
	"net/http"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func Ok(w http.ResponseWriter, data interface{}) {
	respond(w, http.StatusOK, APIResponse{
		Status: "OK",
		Data:   data,
	})
}
func Error(w http.ResponseWriter, code int, message string) {
	respond(w, code, APIResponse{
		Status:  "ERROR",
		Message: message,
	})
}

func respond(w http.ResponseWriter, code int, body APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(body)
}

func DecodeJSONBody[T any](w http.ResponseWriter, r *http.Request) (T, error) {
	var body T
	err := json.NewDecoder(r.Body).Decode(&body)
	return body, err
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
