package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tukangk3tik/rag-starter/internal/chat"
	"github.com/tukangk3tik/rag-starter/internal/config"
	"github.com/tukangk3tik/rag-starter/internal/embedder"
	"github.com/tukangk3tik/rag-starter/internal/prompt"
	"github.com/tukangk3tik/rag-starter/internal/qdrant"
)

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
		config.DefaultConfig.TopK,
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("============== Retrieved Context ==============")

	topKResults := make([]prompt.SearchResult, 0)
	for _, result := range results {
		// debug print the retrieved context
		fmt.Printf(
			"Score: %f\nFile: %s\nChunk: %s\nTitle: %s\nSection: %s\nChunkIndex: %d\n\nContent: \n%s\n",
			result.Score,
			result.Payload.File,
			result.Payload.ChunkID,
			result.Payload.Title,
			result.Payload.Section,
			result.Payload.ChunkIndex,
			result.Payload.Content,
		)
		fmt.Println("--------------------------------------------------")
		if result.Score < config.DefaultConfig.MinScore {
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

	fmt.Print("Response ----------------------------------------\n\n")
	fmt.Print(response)
}
