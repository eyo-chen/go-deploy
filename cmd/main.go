package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

// Response structure
type Response struct {
	Message string `json:"message"`
}

// helloHandler responds with a hello message
func helloHandler(w http.ResponseWriter, r *http.Request) {
	secret := os.Getenv("MY_SECRET")
	port := os.Getenv("PORT")
	response := Response{Message: fmt.Sprintf("Hello, World! %s from %s", secret, port)}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// goodbyeHandler responds with a goodbye message
func goodbyeHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{Message: "Goodbye, World!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/api/hello", helloHandler)
	http.HandleFunc("/api/goodbye", goodbyeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "4040"
	}
	log.Println("Starting server on port", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
