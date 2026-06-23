package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Chat interface {
	Chat(ctx context.Context, messages string) (string, error)
}

// embedding dimensions for Ollama is 768
type OllamaChat struct {
	BaseURL string
}

type chatRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type chatResponse struct {
	Response string `json:"response"`
}

func (c *OllamaChat) Chat(ctx context.Context, messages string) (string, error) {
	reqBody := chatRequest{
		Model:  "qwen3:4b",
		Prompt: messages,
		Stream: false,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshall request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.BaseURL+"/api/generate",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return result.Response, nil
}
