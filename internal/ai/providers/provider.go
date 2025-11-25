package providers

import "context"

// Message represents a chat message
type Message struct {
	Role    string
	Content string
}

// AIProvider defines the interface for AI providers
type AIProvider interface {
	Completion(ctx context.Context, messages []Message) (string, error)
}



