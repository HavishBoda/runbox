package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "strings"
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

// calls the LLM on various messages
func callLLM(messages []Message) (string, error) {
	req := GroqRequest{
		Model: "llama-3.3-70b-versatile",
		Messages: messages,
	}

	jsonReq, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

    httpReq, err := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(jsonReq))
    if err != nil {
        return "", err
    }
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+os.Getenv("GROQ_API_KEY"))

    httpClient := &http.Client{}
    resp, err := httpClient.Do(httpReq)
    if err != nil {
        return "", err
    }

    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    // unmarshal response
    var groqResp GroqResponse
    err = json.Unmarshal(body, &groqResp)
    if err != nil {
        return "", err
    }
    if len(groqResp.Choices) > 0 {
        return groqResp.Choices[0].Message.Content, nil
    }
    return "", fmt.Errorf("no response from LLM")
}

// used to remove ticks from llm response
func extractCode(response string) string {
    response = strings.TrimSpace(response)
    if strings.HasPrefix(response, "```") {
        lines := strings.Split(response, "\n")
        lines = lines[1 : len(lines)-1]
        response = strings.Join(lines, "\n")
    }
    return strings.TrimSpace(response)
}