package prompt

import (
	"fmt"
	"strings"
)

type SearchResult struct {
	Content string
	File    string
	Score   float64
}

func BuildPrompt(
	query string,
	searchResults []SearchResult,
) string {
	var contextBuilder strings.Builder
	for _, result := range searchResults {
		contextBuilder.WriteString(fmt.Sprintf("File: %s\nScore: %f\nContent: %s\n\n", result.File, result.Score, result.Content))
	}
	return fmt.Sprintf(`
		You are a helpful assistant. 
		Use the following context to answer the question.
		If the context does not contain the answer, say "I don't know".

		Context:
		%s

		Question: 
		%s

		Answer: 
	`,
		contextBuilder.String(),
		query,
	)
}
