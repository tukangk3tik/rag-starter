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

var MinScore = float32(0.6)

func main() {
	query := "kenapa redis cepat?"

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

	topKResults := make([]prompt.SearchResult, 0)
	for _, result := range results {
		if result.Score < float64(MinScore) {
			continue
		}
		topKResults = append(topKResults, prompt.SearchResult{
			Content: result.Payload.Content,
			File:    result.Payload.File,
			Score:   result.Score,
		})
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
