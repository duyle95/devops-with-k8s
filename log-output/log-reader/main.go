package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
)

func main() {
	port := os.Getenv("PORT")
	// Port 3000 is important, so as to we can access the service from the outside
	if port == "" {
		port = "3000"
	}

	// Generate a random string (UUID)
	randomString := uuid.New().String()

	// reads /usr/src/app/files/log.txt file and outputs it with a the random string above
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := os.Open("/usr/src/app/files/log.txt")
		if err != nil {
			log.Println("Error opening file:", err)
			http.Error(w, "Error opening file", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		content, err := os.ReadFile("/usr/src/app/files/log.txt")
		if err != nil {
			log.Println("Error reading file:", err)
			http.Error(w, "Error reading file", http.StatusInternalServerError)
			return
		}
		// Write the random string and file content to the response
		w.Write([]byte(fmt.Sprintf("Random String: %s\nFile Content:\n%s", randomString, content)))
	})

	// Start the server and log the port
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
