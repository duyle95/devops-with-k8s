package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	todos := make([]string, 0)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/api/todos", func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(todos)
		if err != nil {
			http.Error(w, "Error marshaling todos", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})

	r.Post("/api/todos", func(w http.ResponseWriter, r *http.Request) {
		// get from form data
		todo := r.FormValue("todo")
		if todo == "" {
			http.Error(w, "Todo cannot be empty", http.StatusBadRequest)
			return
		}
		todos = append(todos, todo)
		w.WriteHeader(http.StatusCreated)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3006"
	}
	fmt.Printf("Starting server on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
