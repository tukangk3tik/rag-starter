package qdrant

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/tukangk3tik/rag-starter/internal/vectordb"
)

type Client struct {
	BaseURL        string
	CollectionName string
}

type SearchResponse struct {
	Result []vectordb.SearchResult `json:"result"`
}

func NewClient(baseUrl string, collection string) *Client {
	return &Client{
		BaseURL:        baseUrl,
		CollectionName: collection,
	}
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

func (c *Client) DeleteCollection(
	ctx context.Context,
) error {
	url := fmt.Sprintf(
		"%s/collections/%s",
		c.BaseURL,
		c.CollectionName,
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		url,
		nil,
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
			"delete collection failed: %d",
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
				"payload": vectordb.PointPayload{
					ChunkID:    point.ID,
					Content:    point.Content,
					File:       point.File,
					Title:      point.Title,
					Section:    point.Section,
					ChunkIndex: point.ChunkIndex,
					IndexedAt:  point.IndexedAt,
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

func (c *Client) Search(ctx context.Context, vector []float32, opts vectordb.SearchOptions) ([]vectordb.SearchResult, error) {
	body := map[string]any{
		"vector":       vector,
		"limit":        opts.TopK,
		"with_payload": true,
	}

	if opts.File != "" {
		body["filter"] = map[string]any{
			"must": []map[string]any{
				{
					"key":   "file",
					"match": map[string]any{"value": opts.File},
				},
			},
		}
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
