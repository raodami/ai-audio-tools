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

// DeepgramClient wraps the Deepgram API for speech recognition
type DeepgramClient struct {
	apiKey string
	client *http.Client
}

// TranscriptionResponse represents the response from Deepgram
type TranscriptionResponse struct {
	Text     string  `json:"text"`
	Duration float64 `json:"duration"`
	Language string  `json:"language"`
}

// NewDeepgramClient creates a new Deepgram client
func NewDeepgramClient() *DeepgramClient {
	return &DeepgramClient{
		apiKey: os.Getenv("DEEPGRAM_API_KEY"),
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

// IsConfigured checks if the API key is set
func (c *DeepgramClient) IsConfigured() bool {
	return c.apiKey != ""
}

// Transcribe sends audio bytes to Deepgram and returns transcription
func (c *DeepgramClient) Transcribe(audio []byte, opts ...TranscribeOption) (*TranscriptionResponse, error) {
	if !c.IsConfigured() {
		return &TranscriptionResponse{
			Text:     "[Deepgram API key not configured - returning placeholder]",
			Duration: float64(len(audio)) / 16000, // rough estimate
			Language: "en",
		}, nil
	}

	params := &TranscribeParams{}
	for _, opt := range opts {
		opt(params)
	}

	// Build request URL
	url := "https://api.deepgram.com/v1/listen"
	if params.Model != "" {
		url += "?model=" + params.Model
	}
	if params.Language != "" {
		url += "&language=" + params.Language
	}

	// Make request
	req, err := http.NewRequest("POST", url, bytes.NewReader(audio))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Token "+c.apiKey)
	req.Header.Set("Content-Type", "audio/wav")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Deepgram API error: %s", string(body))
	}

	// Parse response
	var result struct {
		Result struct {
			Channels []struct {
				Alternatives []struct {
					Text     string  `json:"text"`
					Confidence float64 `json:"confidence"`
				} `json:"alternatives"`
			} `json:"channels"`
			Metadata struct {
				Duration float64 `json:"duration"`
			} `json:"metadata"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Result.Channels) == 0 || len(result.Result.Channels[0].Alternatives) == 0 {
		return nil, fmt.Errorf("no transcription result")
	}

	return &TranscriptionResponse{
		Text:     result.Result.Channels[0].Alternatives[0].Text,
		Duration: result.Result.Metadata.Duration,
		Language: params.Language,
	}, nil
}

// TranscribeFromFile reads audio file and transcribes it
func (c *DeepgramClient) TranscribeFromFile(path string, opts ...TranscribeOption) (*TranscriptionResponse, error) {
	audio, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio file: %w", err)
	}
	return c.Transcribe(audio, opts...)
}

// TranscribeParams holds optional parameters for transcription
type TranscribeParams struct {
	Model    string
	Language string
}

// TranscribeOption is a function that modifies TranscribeParams
type TranscribeOption func(*TranscribeParams)

// WithModel sets the transcription model
func WithModel(model string) TranscribeOption {
	return func(p *TranscribeParams) {
		p.Model = model
	}
}

// WithLanguage sets the expected language
func WithLanguage(lang string) TranscribeOption {
	return func(p *TranscribeParams) {
		p.Language = lang
	}
}

// MockDeepgramClient is used in tests
type MockDeepgramClient struct {
	Response *TranscriptionResponse
	Error    error
}

func (m *MockDeepgramClient) Transcribe(audio []byte, opts ...TranscribeOption) (*TranscriptionResponse, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Response, nil
}
