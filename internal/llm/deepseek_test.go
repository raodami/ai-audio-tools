package llm

import (
	"errors"
	"testing"
)

func TestNewDeepSeekClient(t *testing.T) {
	client := NewDeepSeekClient()
	if client == nil {
		t.Fatal("Expected non-nil client")
	}
}

func TestDeepSeekClient_IsConfigured(t *testing.T) {
	client := NewDeepSeekClient()
	if client.IsConfigured() {
		t.Error("Expected client to not be configured without API key")
	}
}

func TestMockDeepSeekClient_Summarize(t *testing.T) {
	mock := &MockDeepSeekClient{
		SummarizeResponse: "这是一段测试总结。",
	}

	result, err := mock.Summarize("测试文本")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result != "这是一段测试总结。" {
		t.Errorf("Expected '这是一段测试总结。', got %s", result)
	}
}

func TestMockDeepSeekClient_Summarize_Error(t *testing.T) {
	mock := &MockDeepSeekClient{
		Error: errors.New("test error"),
	}

	_, err := mock.Summarize("测试文本")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestMockDeepSeekClient_AnalyzeContract(t *testing.T) {
	mock := &MockDeepSeekClient{
		AnalyzeResponse: &AnalysisResult{
			Summary: "测试合同总结",
			Risks: []RiskItem{
				{Clause: "赔偿条款", RiskLevel: "high", Explanation: "赔偿金额过高"},
			},
			Suggestions: []string{"建议降低赔偿上限"},
		},
	}

	result, err := mock.AnalyzeContract("测试合同文本")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if result.Summary != "测试合同总结" {
		t.Errorf("Expected '测试合同总结', got %s", result.Summary)
	}
	if len(result.Risks) != 1 {
		t.Errorf("Expected 1 risk, got %d", len(result.Risks))
	}
}

func TestMockDeepSeekClient_AnalyzeContract_Error(t *testing.T) {
	mock := &MockDeepSeekClient{
		Error: errors.New("analyze error"),
	}

	_, err := mock.AnalyzeContract("测试合同文本")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}
