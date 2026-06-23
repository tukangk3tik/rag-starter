package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tukangk3tik/rag-starter/internal/chat"
	"github.com/tukangk3tik/rag-starter/internal/embedder"
	"github.com/tukangk3tik/rag-starter/internal/prompt"
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

	topKResults := make([]prompt.SearchResult, len(results))
	for i, result := range results {
		topKResults[i] = prompt.SearchResult{
			Content: result.Payload.Content,
			File:    result.Payload.File,
			Score:   result.Score,
		}
	}

	promptResult := prompt.BuildPrompt(
		query,
		topKResults,
	)

	chatClient := &chat.OllamaChat{
		BaseURL: "http://localhost:11434",
	}

	response, err := chatClient.Chat(
		context.Background(),
		promptResult,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Response:", response)
}
