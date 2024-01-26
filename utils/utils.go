package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
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

func EncodeToJPEG(img image.Image) ([]byte, error) {
	var buffer bytes.Buffer
	err := jpeg.Encode(&buffer, img, nil)
	return buffer.Bytes(), err
}
