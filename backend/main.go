package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	cache = NewLRUCache(1024)

	mux := http.NewServeMux()
	mux.HandleFunc("/get", getHandler)
	mux.HandleFunc("/set", setHandler)

	// Create a CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Allow requests from the React app
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
	})

	// Wrap the server with CORS middleware
	handler := c.Handler(mux)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
