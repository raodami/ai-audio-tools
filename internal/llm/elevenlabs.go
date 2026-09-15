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
	Category  string `json:"category,omitempty"`
}

// CloneRequest represents a voice cloning request
type CloneRequest struct {
	Name          string                 `json:"name"`
	Samples       []SampleFile           `json:"samples,omitempty"`
	Description   string                 `json:"description,omitempty"`
	LabSamples    []LabSample            `json:"lab_samples,omitempty"`
}

// SampleFile represents an audio sample for cloning
type SampleFile struct {
	FileID string `json:"file_id"`
	Labels map[string]string `json:"labels,omitempty"`
}

// LabSample represents a pre-transcribed sample
type LabSample struct {
	Text       string `json:"text"`
	AudioData  []byte `json:"audio_data"`
}

// SynthesizeRequest represents a synthesis request
type SynthesizeRequest struct {
	Text         string            `json:"text"`
	VoiceID      string            `json:"voice_id"`
	ModelID      string            `json:"model_id"`
	VoiceSettings map[string]interface{} `json:"voice_settings,omitempty"`
}

// NewElevenLabsClient creates a new ElevenLabs client
func NewElevenLabsClient() *ElevenLabsClient {
	return &ElevenLabsClient{
		apiKey: os.Getenv("ELEVENLABS_API_KEY"),
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

// IsConfigured checks if the API key is set
func (c *ElevenLabsClient) IsConfigured() bool {
	return c.apiKey != ""
}

// ListVoices returns available voices
func (c *ElevenLabsClient) ListVoices() ([]Voice, error) {
	if !c.IsConfigured() {
		return []Voice{}, fmt.Errorf("ElevenLabs API key not configured")
	}

	req, err := http.NewRequest("GET", "https://api.elevenlabs.io/v1/voices", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
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

// CloneVoice creates a voice clone from audio samples
func (c *ElevenLabsClient) CloneVoice(name, description string, audioFiles []string) (*Voice, error) {
	if !c.IsConfigured() {
		return nil, fmt.Errorf("ElevenLabs API key not configured")
	}

	// Create voice metadata
	cloningSettings := map[string]interface{}{
		"name": name,
		"samples": audioFiles,
	}

	body, err := json.Marshal(cloningSettings)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://api.elevenlabs.io/v1/voice-synthesis/clone", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse response (simplified)
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	return &Voice{
		VoiceID:   result["voice_id"].(string),
		Name:      name,
		Description: description,
	}, nil
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
		return nil, err
	}

	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voiceID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "audio/mpeg")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
