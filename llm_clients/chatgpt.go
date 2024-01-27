package chatgptclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	_ "image/png"
	"io"
	"net/http"
)

type ChatGPTClientInterface interface {
	SendChatCompletionRequest(prompt string) (*ChatCompletionResponse, error)
	SendImageRequest(prompt string) (image.Image, error)
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

type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	SystemFingerprint string   `json:"system_fingerprint"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
}

type ImageRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	n      int    `json:n` // number of images to generate
}

type ImageUrls struct {
	URL string `json:"url"`
}

type ImageResponse struct {
	Created int64       `json:"created"`
	Data    []ImageUrls `json:"data"`
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

func (c *ChatGPTClient) SendImageRequest(prompt string) (image.Image, error) {
	promptFormat := "Generate an image from this description: %s"
	prompt = fmt.Sprintf(promptFormat, prompt)

	requestBody, err := json.Marshal(ImageRequest{
		Model:  "dall-e-2",
		Prompt: prompt,
		n:      1,
	})

	if err != nil {
		fmt.Println("Failed to marshal request body")
		return nil, err
	}

	req, err := http.NewRequest("POST", c.BaseURL+"images/generations", bytes.NewBuffer(requestBody))
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

	var respJson ImageResponse
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		fmt.Println("error parsing to json")
		return nil, err
	}

	// download image
	img, err := downloadImage(respJson.Data[0].URL)
	if err != nil {
		fmt.Println("error downloading image")
		return nil, err
	}

	return img, nil
}

func downloadImage(url string) (image.Image, error) {
	response, err := http.Get(url)
	if err != nil || response.StatusCode != 200 {
		fmt.Println("error getting image from url")
		return nil, err
	}
	defer response.Body.Close()

	img, _, err := image.Decode(response.Body)
	return img, err
}
