package llm

import (
	"errors"
	"testing"
)

func TestNewDeepgramClient(t *testing.T) {
	client := NewDeepgramClient()
	if client == nil {
		t.Fatal("Expected non-nil client")
	}
}

func TestDeepgramClient_IsConfigured(t *testing.T) {
	client := NewDeepgramClient()
	// Without API key, should not be configured
	if client.IsConfigured() {
		t.Error("Expected client to not be configured without API key")
	}
}

func TestMockDeepgramClient_Transcribe_Success(t *testing.T) {
	mock := &MockDeepgramClient{
		Response: &TranscriptionResponse{
			Text:     "Hello world, this is a test.",
			Duration: 5.2,
			Language: "en",
		},
	}

	resp, err := mock.Transcribe([]byte("fake audio"))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if resp.Text != "Hello world, this is a test." {
		t.Errorf("Expected 'Hello world, this is a test.', got %s", resp.Text)
	}
	if resp.Duration != 5.2 {
		t.Errorf("Expected 5.2 duration, got %f", resp.Duration)
	}
}

func TestMockDeepgramClient_Transcribe_Error(t *testing.T) {
	mock := &MockDeepgramClient{
		Error: errors.New("test error"),
	}

	_, err := mock.Transcribe([]byte("fake audio"))
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestTranscribeOption_Model(t *testing.T) {
	params := &TranscribeParams{}
	opt := WithModel("general")
	opt(params)
	if params.Model != "general" {
		t.Errorf("Expected 'general', got %s", params.Model)
	}
}

func TestTranscribeOption_Language(t *testing.T) {
	params := &TranscribeParams{}
	opt := WithLanguage("es")
	opt(params)
	if params.Language != "es" {
		t.Errorf("Expected 'es', got %s", params.Language)
	}
}
