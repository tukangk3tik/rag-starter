package prompt

import (
	"fmt"
	"strings"
)

type SearchResult struct {
	Content string
	File    string
	Score   float32
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

		Use the following language to answer the question based on the question's language.

		Always cite the source files used.
		
		At the end of the answer, list all source filenames. 
		Example: [Source: redis.md]

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
