package config

type RAGConfig struct {
	TopK     int
	MinScore float32
}

var DefaultConfig = RAGConfig{
	TopK:     3,
	MinScore: 0.60,
}
