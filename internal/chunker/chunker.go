package chunker

import (
	"fmt"
	"strings"
)

type Chunk struct {
	ID      string
	File    string
	Content string
	Vector  []float32
}

// need to refactor to split by paragraph, sentence, or other method to make more semantic chunck
func Chunker(text string, filename string) []Chunk {
	var chunks []Chunk

	paragraphs := strings.Split(text, "\n\n")
	for i, p := range paragraphs {
		chunks = append(
			chunks,
			Chunk{
				ID:      fmt.Sprintf("%s-%d", filename, i),
				File:    filename,
				Content: p,
				Vector:  []float32{},
			},
		)
	}

	return chunks
}
