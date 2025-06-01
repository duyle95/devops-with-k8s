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

const todoBackendURL = "http://localhost:8081/api/todos"

// const todoBackendURL = "http://localhost:3006/api/todos"

func main() {
	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "3005"
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

		// // get the list of todos from the todo backend service
		// resp, err := http.Get(todoBackendURL)
		// if err != nil {
		// 	http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		// 	log.Printf("Error fetching todos: %v", err)
		// 	return
		// }
		// defer resp.Body.Close()
		// body, err := io.ReadAll(resp.Body)
		// if err != nil {
		// 	http.Error(w, "Failed to read todos response", http.StatusInternalServerError)
		// 	log.Printf("Error reading todos response: %v", err)
		// 	return
		// }

		// var data []string
		// err = json.Unmarshal(body, &data)
		// if err != nil {
		// 	http.Error(w, "Failed to parse todos response", http.StatusInternalServerError)
		// 	log.Printf("Error parsing todos response: %v", err)
		// 	return
		// }

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
		<form action="%s" method="post">
			<input type="text" name="todo" placeholder="Enter your todo (max 140 characters)" maxlength="140" required>
			<button type="submit">Send</button>
		</form>
		<ul>
			
		`, imageBody, todoBackendURL)
		w.Write([]byte(htmlResponse))
		// for i := 0; i < len(data); i++ {
		// 	if string(data[i]) != "" {
		// 		w.Write([]byte(`
		// 			<li>` + string(data[i]) + `</li>
		// 		`))
		// 	}
		// }
		w.Write([]byte(`
		</ul>
		</body>
		<script>
			// create a script to get todos from the backend service and update the list
			setInterval(() => {
				fetch("` + todoBackendURL + `")
					.then(response => response.json())
					.then(data => {
						const ul = document.querySelector("ul");
						ul.innerHTML = "";
						data.forEach(todo => {
							if (todo) {
								const li = document.createElement("li");
								li.textContent = todo;
								ul.appendChild(li);
							}
						});
					})
					.catch(error => {
						console.error("Error fetching todos:", error);
					});
			}, 5000);
			// Add a listener to the form to prevent default submission and use fetch instead
			document.querySelector("form").addEventListener("submit", function(event) {
				event.preventDefault();
				const formData = new FormData(this);
				fetch(this.action, {
					method: "POST",
					body: formData
				})
				.then(response => {
					if (!response.ok) {
						throw new Error("Network response was not ok");
					}
					return response.text();
				})
			})
		</script>
		</html>
		`))
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
