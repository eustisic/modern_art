package handlers

import (
	"errors"
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

	if q == "" {
		utils.EncodeError(w, "No query string found")
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	// this function will check db for prompt - the prompt will be the name of an artist
	var prompt string
	var found bool
	var err error

	if prompt, found = kv.Search(q); !found {
		// query API and get prompt then insert into db
		prompt, err = GetPrompt(kv, q)
		if err != nil {
			utils.EncodeError(w, err.Error())
			http.Error(w, "Bad Request", http.StatusBadRequest)
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

	if len(resp.Choices) == 0 {
		return "", errors.New("invalid response from chat API")
	}

	return resp.Choices[0].Message.Content, nil
}
