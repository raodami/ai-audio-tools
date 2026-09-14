package processor

import (
	"fmt"
	"os"
	"strings"

	"ai-audio-tools/internal/llm"
)

// AudioProcessor handles audio processing pipeline
type AudioProcessor struct {
	deepgram *llm.DeepgramClient
	deepseek *llm.DeepSeekClient
}

// ProcessResult holds the result of audio processing
type ProcessResult struct {
	Transcript string `json:"transcript"`
	Summary    string `json:"summary"`
}

// NewAudioProcessor creates a new processor with configured clients
func NewAudioProcessor() *AudioProcessor {
	return &AudioProcessor{
		deepgram: llm.NewDeepgramClient(),
		deepseek: llm.NewDeepSeekClient(),
	}
}

// ProcessAudio transcribes audio bytes and summarizes the transcript
func (p *AudioProcessor) ProcessAudio(audioData []byte, fileName string) (*ProcessResult, error) {
	// Step 1: Transcribe with Deepgram
	transcript, err := p.deepgram.Transcribe(audioData)
	if err != nil {
		return nil, fmt.Errorf("transcription failed: %w", err)
	}

	// Extract transcript text from response
	transcriptText := ""
	if transcript != nil && transcript.Text != "" {
		transcriptText = transcript.Text
	}

	// Step 2: Summarize with DeepSeek
	summary, err := p.deepseek.Summarize(transcriptText)
	if err != nil {
		// Return transcript even if summary fails
		return &ProcessResult{
			Transcript: transcriptText,
			Summary:    "[Summary generation failed]",
		}, nil
	}

	return &ProcessResult{
		Transcript: transcriptText,
		Summary:    summary,
	}, nil
}

// ProcessText summarizes text content
func (p *AudioProcessor) ProcessText(text string) (string, error) {
	return p.deepseek.Summarize(text)
}

// IsConfigured checks if all required API keys are set
func (p *AudioProcessor) IsConfigured() bool {
	return p.deepgram.IsConfigured() && p.deepseek.IsConfigured()
}

// GetStatus returns the configuration status
func (p *AudioProcessor) GetStatus() map[string]bool {
	return map[string]bool{
		"deepgram":   p.deepgram.IsConfigured(),
		"deepseek":   p.deepseek.IsConfigured(),
		"elevenlabs": llm.NewElevenLabsClient().IsConfigured(),
	}
}

// LoadTextFromFile reads text content from a file
func LoadTextFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

// FormatTranscript formats transcript for display
func FormatTranscript(result *ProcessResult) string {
	var sb strings.Builder
	sb.WriteString("## Transcript\n\n")
	sb.WriteString(result.Transcript)
	sb.WriteString("\n\n## Summary\n\n")
	sb.WriteString(result.Summary)
	return sb.String()
}
