package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

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

	err = qdrantClient.DeleteCollection(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	err = qdrantClient.CreateCollection(context.Background(), 768)
	if err != nil {
		log.Fatal(err)
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

		chunks := chunker.Chunker(string(content), entry.Name())

		fmt.Println("=====================================")
		fmt.Printf("File: %s\n", filePath)
		// temporary print chunck result
		for _, ch := range chunks {
			vector, err := embedder.Embed(context.Background(), ch.Content)
			if err != nil {
				log.Fatal(err)
			}
			ch.Vector = vector

			timeStamp := time.Now().Format("2006-01-02T15:04:05Z07:00")
			err = qdrantClient.Upsert(
				context.Background(),
				qdrant.Point{
					ID:         ch.ID,
					Vector:     ch.Vector,
					Content:    ch.Content,
					File:       ch.File,
					Title:      ch.Title,
					Section:    ch.Section,
					ChunkIndex: ch.ChunkIndex,
					IndexedAt:  timeStamp,
				},
			)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("-------------------------------------")
			fmt.Printf("ID	   : %s\n", ch.ID)
			fmt.Printf("File	   : %s\n", ch.File)
			fmt.Printf("Content	   : %s\n", ch.Content)
			fmt.Printf("Vector	   : %f\n", ch.Vector[:5])
			fmt.Printf("Title	   : %s\n", ch.Title)
			fmt.Printf("Section	   : %s\n", ch.Section)
			fmt.Printf("ChunkIndex : %d\n", ch.ChunkIndex)
			fmt.Printf("IndexedAt  : %s\n", timeStamp)
			fmt.Printf("Status	   : Success\n")
		}

		fmt.Println()
	}
}
