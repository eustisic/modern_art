package handlers

import (
	"fmt"
	"net/http"

	"modern_art/kvstore"
	chatgptclient "modern_art/llm_clients"
)

func PostPrompt(kv kvstore.StoreInterface, gptClient chatgptclient.ChatGPTClientInterface, w http.ResponseWriter, r *http.Request) {
	// get serarch string from prompt
	query := r.URL.Query()
	q := query.Get("q")

	if q == "" {
		http.Error(w, "Bad Request: No query string found", http.StatusBadRequest)
		return
	}
	// this function will check db for prompt - the prompt will be the name of an artist
	var prompt string
	var found bool
	var err error

	if prompt, found = kv.Search(q); !found {
		// query API and get prompt then insert into db
		prompt, err = GetPrompt(kv, gptClient, q)
		if err != nil {
			http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
			return
		}

		kv.Insert(q, prompt)
	}

	GenerateImage(gptClient, prompt)
}

func PostPromptHandler(kv kvstore.StoreInterface, gptClient chatgptclient.ChatGPTClientInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		PostPrompt(kv, gptClient, w, r)
	}
}

func GetPrompt(kv kvstore.StoreInterface, gptClient chatgptclient.ChatGPTClientInterface, q string) (string, error) {

	resp, err := gptClient.SendChatCompletionRequest(q)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func GenerateImage(gptClient chatgptclient.ChatGPTClientInterface, prompt string) {

	err := gptClient.SendImageRequest(prompt)

	if err != nil {
		fmt.Println(err)
	}
}
