package main

import (
	"fmt"
	"log"
	"modern_art/handlers"
	"modern_art/kvstore"
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
	kv.Insert("Hasui Kawase", "Serene landscapes, traditional woodblock prints, atmospheric, detailed, muted colors, realism, influenced by Western Impressionism")

	router := mux.NewRouter()

	// define info route
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, you are connected to modern art api!")
	})

	// define routes
	router.HandleFunc("/prompt", handlers.PostPromptHandler(kv)).Methods("POST")

	http.Handle("/", router)

	fmt.Println("Server is starting on port: ", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":5000", nil))
}
