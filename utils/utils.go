package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func EncodeMessage(w http.ResponseWriter, message string) {
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func EncodeError(w http.ResponseWriter, err string) {
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}

func LogObject(response interface{}) {
	logStr := fmt.Sprintf("%+v", response)
	fmt.Println("Response object for debugging:", logStr)
}
