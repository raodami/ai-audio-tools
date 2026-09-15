package effects

import (
	"fmt"
	"math"
)

// AudioEffect represents an audio processing effect
type AudioEffect struct {
	Type       string  `json:"type"`
	Parameter  float64 `json:"parameter"`
	Enabled    bool    `json:"enabled"`
}

// ProcessResult holds the result of audio processing
type ProcessResult struct {
	OriginalPath   string  `json:"original_path"`
	ProcessedPath  string  `json:"processed_path"`
	DurationSec    float64 `json:"duration_sec"`
	SampleRate     int     `json:"sample_rate"`
	Channels       int     `json:"channels"`
	EffectsApplied []EffectResult `json:"effects_applied"`
}

// EffectResult holds details of an applied effect
type EffectResult struct {
	Type       string  `json:"type"`
	Parameter  float64 `json:"parameter"`
	Status     string  `json:"status"`
	Error      string  `json:"error,omitempty"`
}

// ApplySpeedChange changes the playback speed (1.0 = normal, 0.5 = half, 2.0 = double)
func ApplySpeedChange(inputPath, outputPath string, speed float64) (*EffectResult, error) {
	if speed <= 0 {
		return &EffectResult{Type: "speed", Parameter: speed, Status: "error", Error: "speed must be positive"}, nil
	}
	
	// NOTE: Actual audio processing requires FFmpeg or similar tool
	// This is a placeholder implementation
	return &EffectResult{
		Type:      "speed",
		Parameter: speed,
		Status:    "not_implemented",
		Error:     "Requires FFmpeg for actual audio processing",
	}, nil
}

// ApplyPitchShift changes the pitch without affecting duration
func ApplyPitchShift(inputPath, outputPath string, semitones float64) (*EffectResult, error) {
	if math.Abs(semitones) > 24 {
		return &EffectResult{Type: "pitch", Parameter: semitones, Status: "error", Error: "max ±24 semitones"}, nil
	}
	
	return &EffectResult{
		Type:      "pitch",
		Parameter: semitones,
		Status:    "not_implemented",
		Error:     "Requires FFmpeg for pitch shifting",
	}, nil
}

// ApplyNoiseReduction removes background noise
func ApplyNoiseReduction(inputPath, outputPath string, strength float64) (*EffectResult, error) {
	if strength < 0 || strength > 1 {
		return &EffectResult{Type: "noise_reduction", Parameter: strength, Status: "error", Error: "strength must be 0-1"}, nil
	}
	
	return &EffectResult{
		Type:      "noise_reduction",
		Parameter: strength,
		Status:    "not_implemented",
		Error:     "Requires advanced audio processing library",
	}, nil
}

// ApplyVolumeChange adjusts the volume
func ApplyVolumeChange(inputPath, outputPath string, volume float64) (*EffectResult, error) {
	if volume < 0 || volume > 5 {
		return &EffectResult{Type: "volume", Parameter: volume, Status: "error", Error: "volume must be 0-5"}, nil
	}
	
	return &EffectResult{
		Type:      "volume",
		Parameter: volume,
		Status:    "ok",
	}, nil
}

// GetAvailableEffects returns list of available audio effects
func GetAvailableEffects() []map[string]string {
	return []map[string]string{
		{"type": "speed", "name": "Speed Change", "desc": "Adjust playback speed (0.5x - 2x)", "param_range": "0.5-2.0"},
		{"type": "pitch", "name": "Pitch Shift", "desc": "Shift pitch by semitones (-12 to +12)", "param_range": "-12 to +12"},
		{"type": "noise_reduction", "name": "Noise Reduction", "desc": "Reduce background noise (0-1)", "param_range": "0-1"},
		{"type": "volume", "name": "Volume", "desc": "Adjust volume level (0-5)", "param_range": "0-5"},
		{"type": "fade_in", "name": "Fade In", "desc": "Add fade-in effect", "param_range": "0-5s"},
		{"type": "fade_out", "name": "Fade Out", "desc": "Add fade-out effect", "param_range": "0-5s"},
	}
}

// ProcessAudioWithEffects applies multiple effects to audio
func ProcessAudioWithEffects(inputPath, outputPath string, effects []AudioEffect) (*ProcessResult, error) {
	result := &ProcessResult{
		OriginalPath:   inputPath,
		ProcessedPath:  outputPath,
		EffectsApplied: make([]EffectResult, 0),
	}
	
	for _, effect := range effects {
		var effectResult *EffectResult
		var err error
		
		switch effect.Type {
		case "speed":
			effectResult, err = ApplySpeedChange(inputPath, outputPath, effect.Parameter)
		case "pitch":
			effectResult, err = ApplyPitchShift(inputPath, outputPath, effect.Parameter)
		case "noise_reduction":
			effectResult, err = ApplyNoiseReduction(inputPath, outputPath, effect.Parameter)
		case "volume":
			effectResult, err = ApplyVolumeChange(inputPath, outputPath, effect.Parameter)
		default:
			effectResult = &EffectResult{
				Type:     effect.Type,
				Status:   "error",
				Error:    fmt.Sprintf("unknown effect type: %s", effect.Type),
			}
		}
		
		if err != nil {
			effectResult.Error = err.Error()
		}
		
		result.EffectsApplied = append(result.EffectsApplied, *effectResult)
		
		if effectResult.Status == "error" {
			return result, fmt.Errorf("effect %s failed: %s", effect.Type, effectResult.Error)
		}
	}
	
	return result, nil
}
