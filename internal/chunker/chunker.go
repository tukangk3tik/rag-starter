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
func Chunker(text string, filename string, size int) []Chunk {
	var chunks []Chunk

	words := strings.Fields(text)
	for i := 0; i < len(words); i += size {
		end := min(i+size, len(words))

		chunks = append(
			chunks,
			Chunk{
				ID:      fmt.Sprintf("%s-%d", filename, i/size),
				File:    filename,
				Content: strings.Join(words[i:end], " "),
				Vector:  []float32{},
			},
		)
	}

	return chunks
}
