package handlers

import (
	"fmt"
	"net/http"
	"os"

	"modern_art/kvstore"
	chatgptclient "modern_art/llm_clients"
	"modern_art/utils"
)

func PostPrompt(kv kvstore.StoreInterface, w http.ResponseWriter, r *http.Request) {
	// get serarch string from prompt
	fmt.Println(r.URL.Query())
	query := r.URL.Query()
	q := query.Get("q")

	fmt.Println(q)

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
		prompt, err = GetPrompt(kv, q)
		if err != nil {
			http.Error(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
			return
		}

		kv.Insert(q, prompt)
	}

	fmt.Println(prompt)

	// if it does not find a promt it will populate the database with a description of that artists style

	// It will then query for art generation

	// The
}

func PostPromptHandler(kv kvstore.StoreInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		PostPrompt(kv, w, r)
	}
}

func GetPrompt(kv kvstore.StoreInterface, q string) (string, error) {
	gptClient := chatgptclient.NewClient(os.Getenv("OPENAI_KEY"), os.Getenv("CHAT_URL"))

	resp, err := gptClient.SendChatCompletionRequest(q)

	if err != nil {
		return "", err
	}

	utils.LogObject(resp)

	return resp.Choices[0].Message.Content, nil
}
