package chatgptclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type ChatGPTClientInterface interface {
	SendRequest(prompt string) (*ResponseBody, error)
}

// ChatGPTClient holds the configuration for the API client
type ChatGPTClient struct {
	APIKey  string
	BaseURL string
}

type OpenAIRequest struct {
	prompt     string
	max_tokens int
}

// NewClient creates a new instance of ChatGPTClient
func NewClient(apiKey, baseURL string) *ChatGPTClient {
	return &ChatGPTClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
}

// RequestBody is the structure of the request body for the ChatGPT API
type RequestBody struct {
	Prompt string `json:"prompt"`
}

// ResponseBody is the structure of the response from the ChatGPT API
type ResponseBody struct {
	Responses []string `json:"responses"`
}

// SendRequest sends a request to the ChatGPT API and returns the response
func (c *ChatGPTClient) SendRequest(prompt string) (*ResponseBody, error) {
	requestBody, err := json.Marshal(RequestBody{Prompt: prompt})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var responseBody ResponseBody
	if err := json.Unmarshal(body, &responseBody); err != nil {
		return nil, err
	}

	return &responseBody, nil
}
