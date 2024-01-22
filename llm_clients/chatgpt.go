package chatgptclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ChatGPTClientInterface interface {
	SendRequest(prompt string) (*ChatCompletionResponse, error)
}

// ChatGPTClient holds the configuration for the API client
type ChatGPTClient struct {
	APIKey  string
	BaseURL string
}

type ChatRequestMessage struct {
	Role    string `json:"role"` // "system", "user", or "assistant"
	Content string `json:"content"`
}

type ChatCompletionRequest struct {
	Model     string               `json:"model"`
	Messages  []ChatRequestMessage `json:"messages"`
	MaxTokens int                  `json:"max_tokens"` // Add this line
}

type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

type Choice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	Logprobs     interface{} `json:"logprobs"` // null or more complex structure
	FinishReason string      `json:"finish_reason"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// NewClient creates a new instance of ChatGPTClient
func NewClient(apiKey, baseURL string) *ChatGPTClient {
	return &ChatGPTClient{
		APIKey:  apiKey,
		BaseURL: baseURL,
	}
}

// SendChatCompletionRequest sends a request to the ChatGPT chat completions endpoint.
func (c *ChatGPTClient) SendChatCompletionRequest(prompt string) (*ChatCompletionResponse, error) {
	content := "In 25 words classify the artistic style of " + prompt

	requestBody, err := json.Marshal(ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []ChatRequestMessage{
			{
				Role:    "user",
				Content: content,
			},
		},
	})

	if err != nil {
		fmt.Println("Failed to marshal request body")
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(requestBody))
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

	var response ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}
