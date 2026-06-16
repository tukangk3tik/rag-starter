package embedder

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Embedder interface {
	Embed(ctx context.Context, text string) ([]float32, error)
}

// embedding dimensions for Ollama is 768
type OllamaEmbedder struct {
	BaseURL string
}

type embeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type embeddingResponse struct {
	Embedding []float32 `json:"embedding"`
}

func (e *OllamaEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	reqBody := embeddingRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshall request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		e.BaseURL+"/api/embeddings",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result embeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result.Embedding, nil
}
