package chatgptclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type ChatGPTClientInterface interface {
	SendChatCompletionRequest(prompt string) (*ChatCompletionResponse, error)
	SendImageRequest(prompt string) error
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
	MaxTokens int                  `json:"max_tokens"`
}

type ImageRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Size   string `json:size`
	n      int    `json:n` // number of images to generate
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
func (c *ChatGPTClient) SendChatCompletionRequest(artist string) (*ChatCompletionResponse, error) {
	promptFormat := "Create a list of 10 comma separated words that describe the style of %s. Limit of 10 words"
	prompt := fmt.Sprintf(promptFormat, artist)

	requestBody, err := json.Marshal(ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []ChatRequestMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		MaxTokens: 50,
	})

	if err != nil {
		fmt.Println("Failed to marshal request body")
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"chat/completions", bytes.NewBuffer(requestBody))
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

func (c *ChatGPTClient) SendImageRequest(prompt string) error {
	promptFormat := "Generate an image from this description: %s"
	prompt = fmt.Sprintf(promptFormat, prompt)

	requestBody, err := json.Marshal(ImageRequest{
		Model:  "dall-e-2",
		Prompt: prompt,
		n:      1,
		Size:   "1024x1024",
	})

	if err != nil {
		fmt.Println("Failed to marshal request body")
		return err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"images/generations", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// upload to S3
	fmt.Println(body)

	return nil
}
