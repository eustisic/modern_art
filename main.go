package main

import (
	"fmt"
	"log"
	"modern_art/handlers"
	"modern_art/kvstore"
	chatgptclient "modern_art/llm_clients"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	kv := kvstore.NewStore()
	gptClient := chatgptclient.NewClient(os.Getenv("OPENAI_KEY"), os.Getenv("OPENAI_URL"))

	router := mux.NewRouter()

	// define info route
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, you are connected to modern art api!")
	})

	// define routes
	router.HandleFunc("/prompt", handlers.PostPromptHandler(kv, gptClient)).Methods("POST")

	http.Handle("/", router)

	fmt.Println("Server is starting on port: ", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":5000", nil))
}
