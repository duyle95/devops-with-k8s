package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
)

// generates a new timestamp every 5 seconds and saves it into a file.
func main() {
	port := os.Getenv("PORT")
	// Port is not needed for this service, but we can set it to 3002
	// This is useful for testing purposes
	// This is internal service that print the timestamp to file
	if port == "" {
		port = "3002"
	}

	// Generate a random string (UUID)
	randomString := uuid.New().String()

	// Print the string every 5 seconds with a timestamp
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var currentString string

	file, err := os.OpenFile("/usr/src/app/files/log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		panic(err)
	}
	defer file.Close()

	go func() {
		for t := range ticker.C {
			currentString = fmt.Sprintf("%s: -- %s\n", t.UTC().Format(time.RFC3339Nano), randomString)
			// Write the string to a file
			if _, err := file.WriteString(currentString); err != nil {
				log.Println("Error writing to file:", err)
			}
			fmt.Println(currentString)
		}
	}()

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		// Temporarily disable reading from file

		// w.Write(fmt.Appendf(nil, "%s\n", currentString))
		// content, err := os.ReadFile("/usr/src/app/files/pingpong.txt")
		// if err != nil {
		// 	log.Println("Error reading file:", err)
		// 	http.Error(w, "Error reading file", http.StatusInternalServerError)
		// 	return
		// }
		// w.Write([]byte(fmt.Sprintf("Ping / Pongs: %s\n", content)))

		w.Write(fmt.Appendf(nil, "%s\n", currentString))
		resp, err := http.Get("http://pingpong-app-svc:2346/pingpong-count")
		if err != nil {
			log.Println("Error fetching pingpong count:", err)
			http.Error(w, "Error fetching pingpong count", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		bodyString := string(bodyBytes)
		w.Write([]byte(fmt.Sprintf("Ping / Pongs: %s\n", bodyString)))
	})

	// Start the server and log the port
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
