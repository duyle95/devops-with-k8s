package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const pathToImage = "/usr/src/app/files/"

// const pathToImage = "./"

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Define a simple handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// get a file from local filesystem and check if it exists and has been created more than 60 minutes ago
		filePath := pathToImage + "random_image.jpg"
		fileInfo, err := os.Stat(filePath)
		var imageBody string
		var fetchImageErr error

		if os.IsNotExist(err) {
			// If the file does not exist, fetch a new image
			imageBody, fetchImageErr = fetchAndStoreNewImage()
			if fetchImageErr != nil {
				http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
				log.Printf("Error fetching image: %v", fetchImageErr)
				return
			}
		} else if fileInfo.ModTime().Add(60 * time.Minute).Before(time.Now()) {
			// If the file is older than 60 minutes, fetch a new image
			imageBody, fetchImageErr = fetchAndStoreNewImage()
			if fetchImageErr != nil {
				http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
				log.Printf("Error fetching image: %v", fetchImageErr)
				return
			}
		} else if fileInfo.ModTime().Add(60 * time.Minute).After(time.Now()) {
			imageBody, fetchImageErr = fetchExistingImage()
			if fetchImageErr != nil {
				http.Error(w, "Failed to read existing image", http.StatusInternalServerError)
				log.Printf("Error reading image: %v", fetchImageErr)
				return
			}
		}

		// Return a simple html response with the image
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		htmlResponse := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">

		</head>
		<body>
		<img src='data:image/jpeg;base64,%s' alt='Random Image' style='max-width:200px;height:200px;'>
		<form action="/submit" method="post">
			<input type="text" name="todo" placeholder="Enter your todo (max 140 characters)" maxlength="140" required>
			<button type="submit">Send</button>
		</form>
		<ul>
			<li>Todo 1</li>
			<li>Todo 2</li>
		</ul>
		</body>
		</html>
		`, imageBody)
		w.Write([]byte(htmlResponse))
	})

	// Start the server and log the port
	log.Printf("Server started on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func fetchAndStoreNewImage() (string, error) {
	// Fetch a random image from Lorem Picsum
	fetchImageURL := "https://picsum.photos/1200"
	resp, err := http.Get(fetchImageURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch image: %s", resp.Status)
	}

	file, err := os.Create(pathToImage + "random_image.jpg")
	if err != nil {
		return "", fmt.Errorf("failed to create image file: %w", err)
	}
	defer file.Close()
	// Copy the image data to the file
	_, err = file.ReadFrom(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save image: %w", err)
	}

	return fetchExistingImage()
}

func fetchExistingImage() (string, error) {
	content, err := os.ReadFile(pathToImage + "random_image.jpg")
	if err != nil {
		return "", fmt.Errorf("failed to read image file: %w", err)
	}
	// Convert the image to a base64 string
	base64Image := base64.StdEncoding.EncodeToString(content)

	return base64Image, nil
}
