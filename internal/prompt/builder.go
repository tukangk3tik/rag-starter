package prompt

import (
	"fmt"
	"strings"

	"github.com/tukangk3tik/rag-starter/internal/vectordb"
)

func BuildPrompt(
	query string,
	searchResults []vectordb.SearchResult,
) string {
	var contextBuilder strings.Builder
	for _, result := range searchResults {
		contextBuilder.WriteString(fmt.Sprintf("File: %s\nScore: %f\nContent: %s\n\n", result.Payload.File, result.Score, result.Payload.Content))
	}
	return fmt.Sprintf(`
		You are a helpful assistant. 
		Use the following context to answer the question.
		If the context does not contain the answer, say "Yo ndak tau kok tanya saya", or "I don't know", depends the question languange.

		Use the following language to answer the question based on the question's language.

		Always cite the source files used. 
		
		At the end of the answer, list all source filenames, except the context does not contain the answer.
		Format: <filename> (<section>) for each source file used. 
		Example below:
		Source: 
		redis.md (Introduction)
		redis.md (Installation)
		redis.md (Configuration)

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
