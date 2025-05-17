package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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
			currentString = fmt.Sprintf("%s\n", t.UTC().Format(time.RFC3339Nano))
			// Write the string to a file
			if _, err := file.WriteString(currentString); err != nil {
				log.Println("Error writing to file:", err)
			}
			fmt.Println(currentString)
		}
	}()

	// Start the server and log the port
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
