package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
)

type Message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

type GroqRequest struct {
    Model    string    `json:"model"`
    Messages []Message `json:"messages"`
}

type GroqResponse struct {
    Choices []struct {
        Message Message `json:"message"`
    } `json:"choices"`
}

func callLLM(messages []Message) (string, error) {
	req := GroqRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: messages,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

}