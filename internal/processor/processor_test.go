package processor

import (
	"testing"
)

func TestNewAudioProcessor(t *testing.T) {
	p := NewAudioProcessor()
	if p == nil {
		t.Fatal("expected non-nil processor")
	}
	if p.deepgram == nil {
		t.Error("expected deepgram client")
	}
	if p.deepseek == nil {
		t.Error("expected deepseek client")
	}
	if p.elevenlabs == nil {
		t.Error("expected elevenlabs client")
	}
}

func TestGetSupportedLanguages(t *testing.T) {
	languages := GetSupportedLanguages()
	if len(languages) == 0 {
		t.Fatal("expected non-empty languages list")
	}

	// Check for key languages
	found := make(map[string]bool)
	for _, lang := range languages {
		found[lang["code"]] = true
	}

	required := []string{"en", "zh", "yue", "ja", "auto"}
	for _, code := range required {
		if !found[code] {
			t.Errorf("expected language code %s", code)
		}
	}
}

func TestGetTranscriptionModels(t *testing.T) {
	models := GetTranscriptionModels()
	if len(models) == 0 {
		t.Fatal("expected non-empty models list")
	}

	found := make(map[string]bool)
	for _, m := range models {
		found[m["code"]] = true
	}

	required := []string{"nova-2", "whisper-large-v3"}
	for _, code := range required {
		if !found[code] {
			t.Errorf("expected model code %s", code)
		}
	}
}

func TestGenerateSubtitleSRT(t *testing.T) {
	transcript := "Hello world this is a test"
	result, err := GenerateSubtitle(transcript, 10.0, "srt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty SRT output")
	}
	if len(result) < 50 {
		t.Error("expected substantial SRT output")
	}
}

func TestGenerateSubtitleVTT(t *testing.T) {
	transcript := "Hello world this is a test"
	result, err := GenerateSubtitle(transcript, 10.0, "vtt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == "" {
		t.Error("expected non-empty VTT output")
	}
	if len(result) < 50 {
		t.Error("expected substantial VTT output")
	}
}

func TestGenerateSubtitleEmpty(t *testing.T) {
	_, err := GenerateSubtitle("", 10.0, "srt")
	if err == nil {
		t.Error("expected error for empty transcript")
	}
}

func TestGenerateSubtitleInvalidFormat(t *testing.T) {
	_, err := GenerateSubtitle("test", 10.0, "pdf")
	if err == nil {
		t.Error("expected error for invalid format")
	}
}
