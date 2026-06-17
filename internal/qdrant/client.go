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
				"payload": map[string]any{
					"chunk_id": point.ID,
					"content":  point.Content,
					"file":     point.File,
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
