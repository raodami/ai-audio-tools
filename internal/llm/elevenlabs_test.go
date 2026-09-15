package llm

import (
	"testing"
)

func TestElevenLabsIsConfigured(t *testing.T) {
	client := NewElevenLabsClient()
	// Without API key, should return false
	if client.IsConfigured() {
		t.Error("Expected IsConfigured to return false when API key is not set")
	}
}

func TestNewElevenLabsClient(t *testing.T) {
	client := NewElevenLabsClient()
	if client == nil {
		t.Fatal("Expected non-nil client")
	}
}

// Note: ListVoices and Synthesize require actual API key
// These tests only verify the client initialization
