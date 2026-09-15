package processor

import (
	"fmt"
	"os"
	"strings"

	"ai-audio-tools/internal/llm"
)

// AudioProcessor handles audio processing pipeline
type AudioProcessor struct {
	deepgram   *llm.DeepgramClient
	deepseek   *llm.DeepSeekClient
	elevenlabs *llm.ElevenLabsClient
}

// IsConfigured checks if the processor has API keys configured
func (p *AudioProcessor) IsConfigured() bool {
	return p.deepgram.IsConfigured() || p.deepseek.IsConfigured() || p.elevenlabs.IsConfigured()
}

// ProcessResult holds the result of audio processing
type ProcessResult struct {
	Transcript string  `json:"transcript"`
	Summary    string  `json:"summary"`
	Language   string  `json:"language"`
	Duration   float64 `json:"duration"`
}

// ProcessOptions holds optional parameters for audio processing
type ProcessOptions struct {
	Language string // "en", "zh", "yue", "ja", etc.
	Model    string // "nova-2", "whisper-large-v3", etc.
}

// NewAudioProcessor creates a new processor with configured clients
func NewAudioProcessor() *AudioProcessor {
	return &AudioProcessor{
		deepgram:   llm.NewDeepgramClient(),
		deepseek:   llm.NewDeepSeekClient(),
		elevenlabs: llm.NewElevenLabsClient(),
	}
}

// ProcessAudio transcribes audio bytes and summarizes the transcript
func (p *AudioProcessor) ProcessAudio(audioData []byte, fileName string, opts *ProcessOptions) (*ProcessResult, error) {
	// Step 1: Transcribe with Deepgram
	var transcript *llm.TranscriptionResponse
	var err error

	if opts != nil && (opts.Language != "" || opts.Model != "") {
		transcript, err = p.deepgram.Transcribe(audioData,
			llm.WithLanguage(opts.Language),
			llm.WithModel(opts.Model),
		)
	} else {
		transcript, err = p.deepgram.Transcribe(audioData)
	}

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
			Summary:    "[Summary unavailable]",
			Language:   transcript.Language,
			Duration:   transcript.Duration,
		}, nil
	}

	return &ProcessResult{
		Transcript: transcriptText,
		Summary:    summary,
		Language:   transcript.Language,
		Duration:   transcript.Duration,
	}, nil
}

// SynthesizeTTS converts text to speech and returns MP3 bytes
func (p *AudioProcessor) SynthesizeTTS(text string, voiceID string) ([]byte, error) {
	if !p.elevenlabs.IsConfigured() {
		return []byte{}, fmt.Errorf("ElevenLabs API key not configured")
	}

	return p.elevenlabs.Synthesize(text, voiceID)
}

// GetVoices returns available voices for TTS
func (p *AudioProcessor) GetVoices() ([]map[string]string, error) {
	voices, err := p.elevenlabs.ListVoices()
	if err != nil {
		return nil, fmt.Errorf("failed to get voices: %w", err)
	}

	result := make([]map[string]string, 0, len(voices))
	for _, v := range voices {
		result = append(result, map[string]string{
			"voice_id":  v.VoiceID,
			"name":      v.Name,
			"description": v.Description,
		})
	}

	return result, nil
}

// GetSupportedLanguages returns list of supported languages
func GetSupportedLanguages() []map[string]string {
	return []map[string]string{
		{"code": "en", "name": "English", "native": "English"},
		{"code": "zh", "name": "Chinese (Mandarin)", "native": "中文"},
		{"code": "yue", "name": "Cantonese", "native": "粤语"},
		{"code": "ja", "name": "Japanese", "native": "日本語"},
		{"code": "ko", "name": "Korean", "native": "한국어"},
		{"code": "fr", "name": "French", "native": "Français"},
		{"code": "de", "name": "German", "native": "Deutsch"},
		{"code": "es", "name": "Spanish", "native": "Español"},
		{"code": "pt", "name": "Portuguese", "native": "Português"},
		{"code": "it", "name": "Italian", "native": "Italiano"},
		{"code": "nl", "name": "Dutch", "native": "Nederlands"},
		{"code": "ar", "name": "Arabic", "native": "العربية"},
		{"code": "hi", "name": "Hindi", "native": "हिन्दీ"},
		{"code": "th", "name": "Thai", "native": "ไทย"},
		{"code": "vi", "name": "Vietnamese", "native": "Tiếng Việt"},
		{"code": "auto", "name": "Auto-detect", "native": "自动检测"},
	}
}

// GetTranscriptionModels returns list of available models
func GetTranscriptionModels() []map[string]string {
	return []map[string]string{
		{"code": "nova-2", "name": "Nova 2 (Recommended)", "desc": "Best quality, supports 30+ languages"},
		{"code": "nova", "name": "Nova", "desc": "Fast and accurate"},
		{"code": "whisper-large-v3", "name": "Whisper Large v3", "desc": "OpenAI Whisper, excellent multilingual"},
		{"code": "whisper-base", "name": "Whisper Base", "desc": "Lightweight and fast"},
	}
}

// GenerateSubtitle generates subtitle content in SRT or VTT format
func GenerateSubtitle(transcript string, duration float64, format string) (string, error) {
	if transcript == "" {
		return "", fmt.Errorf("transcript is empty")
	}

	lines := strings.Fields(transcript)
	if len(lines) == 0 {
		return "", fmt.Errorf("no words in transcript")
	}

	switch format {
	case "srt":
		return generateSRTEntries(transcript, duration), nil
	case "vtt":
		return generateVTTEntries(transcript, duration), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func generateSRTEntries(transcript string, duration float64) string {
	lines := strings.Fields(transcript)
	totalTime := duration
	timePerWord := totalTime / float64(len(lines))

	var sb strings.Builder
	sb.WriteString("SUBTITLES\n\n")

	for i, word := range lines {
		startTime := timePerWord * float64(i)
		endTime := timePerWord * float64(i+1)

		startFormatted := formatTime(startTime)
		endFormatted := formatTime(endTime)

		sb.WriteString(fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, startFormatted, endFormatted, word))
	}

	return sb.String()
}

func generateVTTEntries(transcript string, duration float64) string {
	lines := strings.Fields(transcript)
	totalTime := duration
	timePerWord := totalTime / float64(len(lines))

	var sb strings.Builder
	sb.WriteString("WEBVTT\n\n")

	for i, word := range lines {
		startTime := timePerWord * float64(i)
		endTime := timePerWord * float64(i+1)

		startFormatted := formatTimeWebVTT(startTime)
		endFormatted := formatTimeWebVTT(endTime)

		sb.WriteString(fmt.Sprintf("%s --> %s\n%s\n\n", startFormatted, endFormatted, word))
	}

	return sb.String()
}

func formatTime(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secs := int(seconds) % 60
	millis := int((seconds - float64(int(seconds))) * 1000)

	return fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, secs, millis)
}

func formatTimeWebVTT(seconds float64) string {
	hours := int(seconds) / 3600
	minutes := (int(seconds) % 3600) / 60
	secs := int(seconds) % 60
	millis := int((seconds - float64(int(seconds))) * 1000)

	return fmt.Sprintf("%02d:%02d:%02d.%03d", hours, minutes, secs, millis)
}

// ExportSubtitle exports subtitle to file
func ExportSubtitle(subtitle string, path string, filename string, format string) error {
	ext := ".srt"
	if format == "vtt" {
		ext = ".vtt"
	}

	fullPath := fmt.Sprintf("%s/%s%s", path, filename, ext)

	return os.WriteFile(fullPath, []byte(subtitle), 0644)
}
