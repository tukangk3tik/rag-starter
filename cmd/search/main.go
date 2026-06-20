package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tukangk3tik/rag-starter/internal/embedder"
	"github.com/tukangk3tik/rag-starter/internal/qdrant"
)

func main() {
	query := "deploy preprod"

	embedder := &embedder.OllamaEmbedder{
		BaseURL: "http://localhost:11434",
	}

	queryVector, err := embedder.Embed(
		context.Background(),
		query,
	)
	if err != nil {
		panic(err)
	}

	qdrantClient := &qdrant.Client{
		BaseURL:        "http://localhost:6333",
		CollectionName: "knowledge",
	}

	results, err := qdrantClient.Search(
		context.Background(),
		queryVector,
		5,
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results {
		fmt.Printf("ID: %s, Content: %s\n", result.ID, result.Payload.Content)
	}
}
