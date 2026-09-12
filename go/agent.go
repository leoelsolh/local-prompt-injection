package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"github.com/leoelsolh/local-prompt-injection/conf"
)

var toolPattern = regexp.MustCompile(`<tool>http_get:([^<]+)</tool>`)

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

func chat(model string, messages []Message) (string, error) {
	req := ChatRequest{
		Model:    model,
		Stream:   false,
		Messages: messages,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(conf.OllamaURL, "application/json", bytes.NewReader(data))
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

func runAgent(model string, prompt string) (string, error) {
	messages := []Message{
		{Role: "system", Content: conf.SystemPrompt},
		{Role: "user", Content: prompt},
	}

	content, err := chat(model, messages)
	if err != nil {
		return "", err
	}

	match := toolPattern.FindStringSubmatch(content)
	if match == nil {
		return content, nil
	}
	targetURL := match[1]

	fetchResp, err := http.Get(targetURL)
	if err != nil {
		return "", err
	}
	defer fetchResp.Body.Close()

	fetchBytes, err := io.ReadAll(fetchResp.Body)
	if err != nil {
		return "", err
	}
	fetched := string(fetchBytes)
	if len(fetched) > conf.MaxContent {
		fetched = fetched[:conf.MaxContent]
	}

	messages = append(messages, Message{Role: "assistant", Content: content})
	messages = append(messages, Message{Role: "user", Content: "Fetched With Tool:\n" + fetched})

	final, err := chat(model, messages)
	if err != nil {
		return "", err
	}
	return final, nil
}
