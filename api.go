package main

import (
	"net/http"
	"encoding/json"
	// "os"
	// "fmt"

	"github.com/joho/godotenv"
)

type UserRequest struct {
	Task string `json:"task"`
}

type UserResponse struct {
	Result string `json:"result"`
	Status string `json:"status"`
}

func handleTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var incoming UserRequest
	err := json.NewDecoder(r.Body).Decode(&incoming)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(UserResponse{Status: "error", Result: "invalid request"})
        return
	}

	output, err := runAgent(incoming.Task)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(UserResponse{Status: "error", Result: err.Error()})
        return
    }

	json.NewEncoder(w).Encode(UserResponse{Status: "success", Result: output})
}

func main() {
    godotenv.Load()

	http.HandleFunc("/v1/task", handleTask)
	http.ListenAndServe(":8080", nil)
}