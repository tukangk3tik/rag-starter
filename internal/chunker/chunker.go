package chunker

import (
	"fmt"
	"strings"
)

type Chunk struct {
	ID         string
	File       string
	Title      string
	Section    string
	ChunkIndex int

	Content string
	Vector  []float32
}

func Chunker(text string, filename string) []Chunk {
	var chunks []Chunk

	paragraphs := strings.Split(text, "\n\n")
	title := ""
	section := ""
	for i, p := range paragraphs {
		if strings.HasPrefix(p, "# ") {
			title = strings.TrimPrefix(p, "# ")
			i--
			continue
		}

		if strings.HasPrefix(p, "## ") {
			section = strings.TrimPrefix(p, "## ")
			i--
			continue
		}

		if strings.TrimSpace(p) == "" {
			i--
			continue
		}
		chunks = append(
			chunks,
			Chunk{
				ID:         fmt.Sprintf("%s-%d", filename, i),
				File:       filename,
				Title:      title,
				Section:    section,
				ChunkIndex: i,
				Content:    p,
				Vector:     []float32{},
			},
		)
	}

	return chunks
}
