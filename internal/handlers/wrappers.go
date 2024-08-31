package handlers

import (
	"encoding/json"
	"net/http"
)

func writeJSONError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	errorResponse := Response{
		Error: message,
	}
	json.NewEncoder(w).Encode(errorResponse)
}
