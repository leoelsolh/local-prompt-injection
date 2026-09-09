package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/leoelsolh/local-prompt-injection/conf"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}
type ChatResponse struct {
	Message Message `json:"message"`
	Done    bool    `json:"done"`
}

func runAgent(model string, prompt string) (string, error) {
	req := ChatRequest{
		Model:  model,
		Stream: false,
		Messages: []Message{
			{Role: "system", Content: conf.SystemPrompt},
			{Role: "user", Content: prompt},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	body := bytes.NewReader(data)
	resp, err := http.Post(conf.OllamaURL, "application/json", body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var cr ChatResponse
	err = json.Unmarshal(respBytes, &cr)
	if err != nil {
		return "", err
	}

	return cr.Message.Content, nil
}
