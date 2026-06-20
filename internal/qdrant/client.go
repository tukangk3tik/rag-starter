package qdrant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Client struct {
	BaseURL        string
	CollectionName string
}

type ChunkPayload struct {
	ChunkID string `json:"chunk_id"`
	Content string `json:"content"`
	File    string `json:"file"`
}

type SearchResult struct {
	ID      string       `json:"id"`
	Version int          `json:"version"`
	Score   float64      `json:"score"`
	Payload ChunkPayload `json:"payload"`
}

type SearchResponse struct {
	Result []SearchResult `json:"result"`
}

func (c *Client) CreateCollection(
	ctx context.Context,
	vectorSize int,
) error {
	body := map[string]any{
		"vectors": map[string]any{
			"size":     vectorSize,
			"distance": "Cosine",
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(
		"%s/collections/%s",
		c.BaseURL,
		c.CollectionName,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		url,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode > 300 {
		return fmt.Errorf(
			"create collection failed: %d",
			resp.StatusCode,
		)
	}

	return nil
}

func (c *Client) Upsert(
	ctx context.Context,
	point Point,
) error {
	body := map[string]any{
		"points": []map[string]any{
			{
				"id":     uuid.New().String(),
				"vector": point.Vector,
				"payload": ChunkPayload{
					ChunkID: point.ID,
					Content: point.Content,
					File:    point.File,
				},
			},
		},
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// create the URL for the upsert request
	url := fmt.Sprintf(
		"%s/collections/%s/points",
		c.BaseURL,
		c.CollectionName,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPut,
		url,
		bytes.NewBuffer(payload),
	)

	if err != nil {
		return err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf(
			"upsert failed: %d",
			resp.StatusCode,
		)
	}

	return nil
}

func (c *Client) Search(
	ctx context.Context,
	vector []float32,
	limit int,
) ([]SearchResult, error) {
	body := map[string]any{
		"vector":       vector,
		"limit":        limit,
		"with_payload": true,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf(
		"%s/collections/%s/points/search",
		c.BaseURL,
		c.CollectionName,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		url,
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf(
			"search failed: %d",
			resp.StatusCode,
		)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return response.Result, nil
}
