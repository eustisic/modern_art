package utils

import (
	"encoding/json"
	"net/http"
)

func EncodeMessage(w http.ResponseWriter, message string) {
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func EncodeError(w http.ResponseWriter, err string) {
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}
