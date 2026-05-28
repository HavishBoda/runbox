package main

import (
	"net/http"
	"encoding/json"
	"fmt"

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

	output, err := runAgent(incoming.Task, func(s string){})
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(UserResponse{Status: "error", Result: err.Error()})
        return
    }

	json.NewEncoder(w).Encode(UserResponse{Status: "success", Result: output})
}

// outputs at each llm step instead of just at the end
func handleTaskStream(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/event-stream")
    w.Header().Set("Cache-Control", "no-cache")
    w.Header().Set("Connection", "keep-alive")

	var incoming UserRequest
	err := json.NewDecoder(r.Body).Decode(&incoming)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(UserResponse{Status: "error", Result: "invalid request"})
        return
	}

	// callback that writes to stream
	stream := func(line string) {
		fmt.Fprintf(w, "data: %s\n\n", line)
		f, ok := w.(http.Flusher)
		if ok {
			f.Flush()
		}
	}

	output, err := runAgent(incoming.Task, stream)
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
	http.HandleFunc("/v1/stream/task", handleTaskStream)
	http.ListenAndServe(":8080", nil)
}