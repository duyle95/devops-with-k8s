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
	var count int64 = -1

	http.HandleFunc("/pingpong", func(w http.ResponseWriter, r *http.Request) {
		count = count + 1
		w.Write([]byte("pong " + strconv.FormatInt(count, 10)))
	})

	// Start the server and log the port
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
