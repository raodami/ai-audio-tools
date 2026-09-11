package llm

import (
	"errors"
	"testing"
)

func TestNewElevenLabsClient(t *testing.T) {
	client := NewElevenLabsClient()
	if client == nil {
		t.Fatal("Expected non-nil client")
	}
}

func TestElevenLabsClient_IsConfigured(t *testing.T) {
	client := NewElevenLabsClient()
	if client.IsConfigured() {
		t.Error("Expected client to not be configured without API key")
	}
}

func TestMockElevenLabsClient_Synthesize(t *testing.T) {
	mock := &MockElevenLabsClient{}

	audio, err := mock.Synthesize("Hello world", "test-voice")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(audio) == 0 {
		t.Error("Expected non-empty audio data")
	}
}

func TestMockElevenLabsClient_Synthesize_Error(t *testing.T) {
	mock := &MockElevenLabsClient{
		Error: errors.New("synthesis failed"),
	}

	_, err := mock.Synthesize("Hello world", "test-voice")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestMockElevenLabsClient_ListVoices(t *testing.T) {
	mock := &MockElevenLabsClient{}

	voices, err := mock.ListVoices()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(voices) != 2 {
		t.Errorf("Expected 2 voices, got %d", len(voices))
	}
}
