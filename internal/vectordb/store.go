package vectordb

import "context"

type SearchOptions struct {
	TopK     int
	MinScore float32
	File     string
	Section  string
}

type PointPayload struct {
	ChunkID    string `json:"chunk_id"`
	Content    string `json:"content"`
	File       string `json:"file"`
	Title      string `json:"title"`
	Section    string `json:"section"`
	ChunkIndex int    `json:"chunk_index"`
	IndexedAt  string `json:"indexed_at"`
}

type SearchResult struct {
	ID      string       `json:"id"`
	Version int          `json:"version"`
	Score   float32      `json:"score"`
	Payload PointPayload `json:"payload"`
}

type VectorStore interface {
	Search(ctx context.Context, vector []float32, opts SearchOptions) ([]SearchResult, error)
}
