package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tukangk3tik/rag-starter/internal/chat"
	"github.com/tukangk3tik/rag-starter/internal/embedder"
	"github.com/tukangk3tik/rag-starter/internal/prompt"
	"github.com/tukangk3tik/rag-starter/internal/qdrant"
	"github.com/tukangk3tik/rag-starter/internal/retriever"
	"github.com/tukangk3tik/rag-starter/internal/vectordb"
)

func main() {
	query := "kenapa redis cepat?"

	embedder := &embedder.OllamaEmbedder{
		BaseURL: "http://localhost:11434",
	}

	qdrantClient := qdrant.NewClient(
		"http://localhost:6333",
		"knowledge",
	)

	config := vectordb.SearchOptions{
		TopK:     5,
		MinScore: 0.6,
		File:     "deployment.md",
	}

	startTime := time.Now()
	re := retriever.NewRetriever(
		embedder,
		qdrantClient,
		config,
	)

	results, err := re.Retrieve(
		context.Background(),
		query,
	)
	if err != nil {
		log.Fatal(err)
	}

	promptResult := prompt.BuildPrompt(
		query,
		results,
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
	fmt.Println(response)

	fmt.Printf("Execution Time: %v\n", time.Since(startTime))
}
