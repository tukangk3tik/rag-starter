# RAG Starter

A small Go-based starter project for building a retrieval-augmented generation (RAG) workflow with Qdrant and Ollama.

## What this project does

The project demonstrates a simple end-to-end RAG pipeline:

1. Read markdown/text files from the `docs/` folder.
2. Split the content into heading-aware chunks and preserve metadata such as title, section, chunk index, and ingestion time.
3. Generate embeddings with Ollama.
4. Store the chunks and their metadata in a Qdrant collection.
5. Search relevant chunks by vector similarity, filter by score threshold, and display the retrieved context for debugging.
6. Build a prompt from the retrieved context, require source citations, and generate a response with Ollama.

## Project structure

- `cmd/ingest/` - ingestion entry point that reads documents and stores embeddings in Qdrant
- `cmd/search/` - search entry point that retrieves relevant chunks and generates a response
- `internal/chunker/` - chunking logic
- `internal/embedder/` - Ollama embedding client
- `internal/chat/` - Ollama chat generation client
- `internal/prompt/` - prompt construction for the LLM
- `internal/qdrant/` - Qdrant client for collection creation, upsert, and search
- `internal/config/` - retrieval configuration defaults for Top-K and score filtering
- `docs/` - sample documents used for ingestion

## Prerequisites

- Go 1.24+
- Docker
- Ollama installed and running locally
- Qdrant running locally (or via Docker)

## Quick start

### 1. Start Qdrant

```bash
docker compose up -d
```

### 2. Pull the required Ollama models

```bash
ollama pull nomic-embed-text
ollama pull qwen3:4b
```

### 3. Run the ingestion flow

```bash
go run ./cmd/ingest
```

This reads files from `docs/`, chunks them, generates embeddings, and stores them in Qdrant.

### 4. Run the search flow

```bash
go run ./cmd/search
```

This embeds a sample query, searches Qdrant for relevant context, filters results by score, builds a prompt with cited sources, and prints the generated answer.

## Configuration

The code currently uses these local defaults:

- Qdrant base URL: `http://localhost:6333`
- Ollama base URL: `http://localhost:11434`
- Qdrant collection name: `knowledge`
- Retrieval Top-K: `3`
- Minimum retrieval score: `0.60`

You can adjust these values directly in the relevant Go files or in the retrieval config package if needed.

## Notes

This repository is a starter implementation and is intended for learning and experimentation with RAG concepts in Go.
