package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tukangk3tik/rag-starter/internal/chunker"
	"github.com/tukangk3tik/rag-starter/internal/embedder"
	"github.com/tukangk3tik/rag-starter/internal/qdrant"
)

func main() {
	dirPath := "./docs"

	qdrantClient := &qdrant.Client{
		BaseURL:        "http://localhost:6333",
		CollectionName: "knowledge",
	}

	embedder := &embedder.OllamaEmbedder{
		BaseURL: "http://localhost:11434",
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())

		content, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal(err)
		}

		chunck := chunker.Chunker(string(content), entry.Name(), 5)

		fmt.Println("=====================================")
		fmt.Printf("File: %s\n", filePath)
		// temporary print chunck result
		for _, ch := range chunck {
			vector, err := embedder.Embed(context.Background(), ch.Content)
			if err != nil {
				log.Fatal(err)
			}
			ch.Vector = vector

			fmt.Println("-------------------------------------")
			fmt.Printf("ID	   : %s\n", ch.ID)
			fmt.Printf("File	   : %s\n", ch.File)
			fmt.Printf("Content	   : %s\n", ch.Content)
			fmt.Printf("Vector	   : %f\n", ch.Vector[:5])

			err = qdrantClient.Upsert(
				context.Background(),
				qdrant.Point{
					ID:      ch.ID,
					Vector:  ch.Vector,
					Content: ch.Content,
					File:    ch.File,
				},
			)
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println()
	}
}
