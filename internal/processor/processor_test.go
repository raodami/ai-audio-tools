package processor

import (
	"strings"
	"testing"
)

func TestNewAudioProcessor(t *testing.T) {
	p := NewAudioProcessor()
	if p == nil {
		t.Fatal("Expected non-nil processor")
	}
}

func TestGetStatus(t *testing.T) {
	p := NewAudioProcessor()
	status := p.GetStatus()
	if _, ok := status["deepgram"]; !ok {
		t.Error("Expected deepgram key in status")
	}
	if _, ok := status["deepseek"]; !ok {
		t.Error("Expected deepseek key in status")
	}
}

func TestFormatTranscript(t *testing.T) {
	result := &ProcessResult{
		Transcript: "Hello world",
		Summary:    "A greeting",
	}
	formatted := FormatTranscript(result)
	if !strings.Contains(formatted, "Hello world") {
		t.Error("Expected transcript in formatted output")
	}
	if !strings.Contains(formatted, "A greeting") {
		t.Error("Expected summary in formatted output")
	}
}
