package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"modern_art/kvstore"
)

func PostPrompt(kv kvstore.StoreInterface, w http.ResponseWriter, r *http.Request) {
	// get serarch string from prompt
	query := r.URL.Query()
	q := query.Get("q")

	if q == "" {
		json.NewEncoder(w).Encode(map[string]string{"message": "No query string found"})
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
	// this function will check db for prompt - the prompt will be the name of an artist
	var prompt string
	var found bool

	if prompt, found = kv.Search(q); !found {
		// query API and get prompt then insert into db
		prompt = GetPrompt(kv, q)
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

func GetPrompt(kv kvstore.StoreInterface, q string) string {
	return ""
}
