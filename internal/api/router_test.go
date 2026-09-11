package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"ai-audio-tools/internal/store"
)

func TestRouterSetup(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s, _ := store.NewDB(":memory:")
	defer s.Close()

	r := gin.New()
	SetupRoutes(r, s)

	// Test GET /api/user/usage
	req, _ := http.NewRequest("GET", "/api/user/usage", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["free_quota"] != float64(30) {
		t.Errorf("Expected free_quota 30, got %v", resp["free_quota"])
	}
}
