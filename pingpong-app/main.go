package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	var count int64 = 0

	http.HandleFunc("/pingpong", func(w http.ResponseWriter, r *http.Request) {
		count = count + 1

		file, err := os.OpenFile("/usr/src/app/files/pingpong.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println("Error opening file:", err)
			http.Error(w, "Error opening file", http.StatusInternalServerError)
			return
		}
		defer file.Close()
		// Write the count to the file
		if _, err := file.WriteString(strconv.FormatInt(count, 10)); err != nil {
			log.Println("Error writing to file:", err)
			http.Error(w, "Error writing to file", http.StatusInternalServerError)
			return
		}

		w.Write([]byte("pong " + strconv.FormatInt(count, 10)))
	})

	// Start the server and log the port
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
