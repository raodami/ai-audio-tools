package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// ElevenLabsClient wraps the ElevenLabs API for text-to-speech
type ElevenLabsClient struct {
	apiKey string
	client *http.Client
}

// Voice represents a voice option
type Voice struct {
	VoiceID   string `json:"voice_id"`
	Name      string `json:"name"`
	Description string `json:"description"`
}

// SynthesizeRequest represents a synthesis request
type SynthesizeRequest struct {
	Text         string            `json:"text"`
	VoiceID      string            `json:"voice_id"`
	ModelID      string            `json:"model_id"`
	VoiceSettings map[string]interface{} `json:"voice_settings"`
}

// NewElevenLabsClient creates a new ElevenLabs client
func NewElevenLabsClient() *ElevenLabsClient {
	return &ElevenLabsClient{
		apiKey: os.Getenv("ELEVENLABS_API_KEY"),
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// IsConfigured checks if the API key is set
func (c *ElevenLabsClient) IsConfigured() bool {
	return c.apiKey != ""
}

// Synthesize converts text to speech and returns MP3 bytes
func (c *ElevenLabsClient) Synthesize(text string, voiceID string) ([]byte, error) {
	if !c.IsConfigured() {
		return []byte(""), nil
	}

	if voiceID == "" {
		voiceID = "pNInz6obpgDQGcFmaJgB" // Adam voice default
	}

	reqBody := SynthesizeRequest{
		Text:      text,
		VoiceID:   voiceID,
		ModelID:   "eleven_monolingual_v1",
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voiceID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("xi-api-key", c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	audio, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ElevenLabs API error: %s", string(audio))
	}

	return audio, nil
}

// ListVoices returns available voices
func (c *ElevenLabsClient) ListVoices() ([]Voice, error) {
	if !c.IsConfigured() {
		return nil, nil
	}

	req, err := http.NewRequest("GET", "https://api.elevenlabs.io/v1/voices", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("xi-api-key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Voices []Voice `json:"voices"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Voices, nil
}

// MockElevenLabsClient is used in tests
type MockElevenLabsClient struct {
	SynthesizeResponse []byte
	Error              error
}

func (m *MockElevenLabsClient) Synthesize(text string, voiceID string) ([]byte, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return []byte("fake-mpeg-audio-data"), nil
}

func (m *MockElevenLabsClient) ListVoices() ([]Voice, error) {
	return []Voice{
		{VoiceID: "test1", Name: "Test Voice 1", Description: "A test voice"},
		{VoiceID: "test2", Name: "Test Voice 2", Description: "Another test voice"},
	}, nil
}
