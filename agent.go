package main

import (
	"fmt"
	"strings"
)

func runAgent(task string) (string, error) {
	prompt := "You are an agent that solves tasks by writing Python code or answering factual questions. " +
    "When you need to write code, wrap it in <code>, </code> tags. " +
    "When the task is complete, respond with <done> followed by the final answer in plain text. " +
    "Never put code after <done>, only the final answer."

	messages := []Message{
		{Role: "system", Content: prompt},
		{Role: "user", Content: task},
	}

	// max 5 iterations before we stop the LLM
	for i := 0; i < 5; i++ {
		response, err := callLLM(messages)
		if err != nil {
			return "", err
		}

		// debugging, checking for <code> tags
		//fmt.Println("LLM response:", response)

		// add LLM response to messages list
		messages = append(messages, Message{Role: "assistant", Content: response})

		// check for <done>, return answer
		if strings.Contains(response, "<done>") {
			text := strings.Split(response, "<done>")
			return strings.TrimSpace(text[1]), nil
		}

		// if llm generates code, run code and hand it back to llm
		if strings.Contains(response, "<code>") {
			text := strings.Split(response, "<code>")
			llmCode := strings.Split(text[1], "</code>")[0]

			output, err := runCode(llmCode)
			if err != nil {
				return "", err
			}
			messages = append(messages, Message{Role: "user", Content: output})
		}
	}
	return "", fmt.Errorf("max iterations reached\n")
}