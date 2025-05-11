package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Generate a random string (UUID)
	randomString := uuid.New().String()

	// Print the string every 5 seconds with a timestamp
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var currentString string

	go func() {
		for t := range ticker.C {
			currentString = fmt.Sprintf("%s: -- %s\n", t.UTC().Format(time.RFC3339Nano), randomString)
			fmt.Println(currentString)
		}
	}()

	// New endpoint
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(currentString))
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// Start the server and log the port
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
