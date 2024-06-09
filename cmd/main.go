package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response structure
type Response struct {
	Message string `json:"message"`
}

// helloHandler responds with a hello message
func helloHandler(w http.ResponseWriter, r *http.Request) {
	a := 1
	response := Response{Message: "Hello, World!"}
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

	log.Println("Starting server on :4040...")
	if err := http.ListenAndServe(":4040", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err.Error())
	}
}
