package qdrant

type Point struct {
	ID         string
	Vector     []float32
	Content    string
	File       string
	Title      string
	Section    string
	ChunkIndex int
	IndexedAt  string
}
