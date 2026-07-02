package retriever

import (
	"context"
	"fmt"
	"log"

	"github.com/tukangk3tik/rag-starter/internal/embedder"
	"github.com/tukangk3tik/rag-starter/internal/vectordb"
)

type Retriever struct {
	Embedder embedder.Embedder
	Store    vectordb.VectorStore
	Config   vectordb.SearchOptions
}

func NewRetriever(
	embedder embedder.Embedder,
	store vectordb.VectorStore,
	config vectordb.SearchOptions,
) *Retriever {
	return &Retriever{
		Embedder: embedder,
		Store:    store,
		Config:   config,
	}
}

func (r *Retriever) Retrieve(
	ctx context.Context,
	query string,
) ([]vectordb.SearchResult, error) {

	// embed process
	queryVector, err := r.Embedder.Embed(
		ctx,
		query,
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	// search in vector store
	results, err := r.Store.Search(
		ctx,
		queryVector,
		r.Config,
	)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	fmt.Println("============== Retrieved Context ==============")

	// filter by score
	topKResults := make([]vectordb.SearchResult, 0)
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
		if result.Score < r.Config.MinScore {
			continue
		}
		topKResults = append(topKResults, vectordb.SearchResult{
			ID:      result.ID,
			Version: result.Version,
			Score:   result.Score,
			Payload: result.Payload,
		})
	}

	return topKResults, nil
}
