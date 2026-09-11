package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// DeepSeekClient wraps the DeepSeek API for LLM calls
type DeepSeekClient struct {
	apiKey string
	baseURL string
	client *http.Client
}

// SummarizeRequest represents a summarize request
type SummarizeRequest struct {
	Model    string `json:"model"`
	Messages []Message `json:"messages"`
}

// SummarizeResponse represents a summarize response
type SummarizeResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// RiskItem represents a risk item in contract analysis
type RiskItem struct {
	Clause      string `json:"clause"`
	RiskLevel   string `json:"risk_level"` // high/medium/low
	Explanation string `json:"explanation"`
}

// AnalysisResult represents the result of contract analysis
type AnalysisResult struct {
	Summary     string      `json:"summary"`
	Risks       []RiskItem  `json:"risks"`
	Suggestions []string    `json:"suggestions"`
}

// NewDeepSeekClient creates a new DeepSeek client
func NewDeepSeekClient() *DeepSeekClient {
	return &DeepSeekClient{
		apiKey:    os.Getenv("DEEPSEEK_API_KEY"),
		baseURL:   "https://api.deepseek.com/v1",
		client:    &http.Client{Timeout: 60 * time.Second},
	}
}

// IsConfigured checks if the API key is set
func (c *DeepSeekClient) IsConfigured() bool {
	return c.apiKey != ""
}

// Summarize sends text to DeepSeek and returns a summary
func (c *DeepSeekClient) Summarize(text string) (string, error) {
	if !c.IsConfigured() {
		return "[DeepSeek API key not configured - returning placeholder]", nil
	}

	msg := Message{
		Role:    "user",
		Content: fmt.Sprintf("请总结以下文本的主要内容，控制在100字以内：\n%s", text),
	}

	reqBody := SummarizeRequest{
		Model:    "deepseek-chat",
		Messages: []Message{msg},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("DeepSeek API error: %s", string(respBody))
	}

	var result SummarizeResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty response from DeepSeek")
	}

	return result.Choices[0].Message.Content, nil
}

// AnalyzeContract analyzes a contract and returns risks and suggestions
func (c *DeepSeekClient) AnalyzeContract(text string) (*AnalysisResult, error) {
	if !c.IsConfigured() {
		return &AnalysisResult{
			Summary:     "[DeepSeek API key not configured - placeholder result]",
			Risks:       []RiskItem{{Clause: "测试条款", RiskLevel: "low", Explanation: "API未配置"}},
			Suggestions: []string{"配置DEEPSEEK_API_KEY环境变量"},
		}, nil
	}

	prompt := `你是一位专业的法律合同审查专家。请分析以下合同文本，识别其中的风险条款，并提供修改建议。

请以JSON格式返回结果，包含以下字段：
- summary: 合同简要概述（100字以内）
- risks: 风险列表，每项包含 clause(条款内容), risk_level(高/中/低), explanation(解释)
- suggestions: 修改建议列表

合同文本：
%s`

	msg := Message{
		Role:    "user",
		Content: fmt.Sprintf(prompt, text),
	}

	reqBody := SummarizeRequest{
		Model:    "deepseek-chat",
		Messages: []Message{msg},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DeepSeek API error: %s", string(respBody))
	}

	var result SummarizeResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("empty response from DeepSeek")
	}

	// Parse JSON result from assistant message
	var analysis AnalysisResult
	assistantMsg := result.Choices[0].Message.Content
	// Handle cases where the model wraps JSON in markdown code blocks
	assistantMsg = strings.TrimPrefix(assistantMsg, "```json")
	assistantMsg = strings.TrimSuffix(assistantMsg, "```")
	assistantMsg = strings.TrimSpace(assistantMsg)

	if err := json.Unmarshal([]byte(assistantMsg), &analysis); err != nil {
		// If parsing fails, return a basic result
		return &AnalysisResult{
			Summary:     assistantMsg[:min(len(assistantMsg), 200)],
			Risks:       []RiskItem{},
			Suggestions: []string{},
		}, nil
	}

	return &analysis, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// MockDeepSeekClient is used in tests
type MockDeepSeekClient struct {
	SummarizeResponse string
	AnalyzeResponse   *AnalysisResult
	Error             error
}

func (m *MockDeepSeekClient) Summarize(text string) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}
	return m.SummarizeResponse, nil
}

func (m *MockDeepSeekClient) AnalyzeContract(text string) (*AnalysisResult, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.AnalyzeResponse, nil
}
