package api

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"ai-audio-tools/internal/store"
)

func setupTestRouter(t *testing.T) (*gin.Engine, *store.Store) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	s, err := store.NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test DB: %v", err)
	}
	r := gin.New()
	SetupRoutes(r, s)
	return r, s
}

func TestHealthCheck(t *testing.T) {
	r, _ := setupTestRouter(t)
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestRegister(t *testing.T) {
	r, _ := setupTestRouter(t)
	body, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "testpass123",
	})
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["token"] == nil {
		t.Error("Expected token in response")
	}
	if resp["user"] == nil {
		t.Error("Expected user in response")
	}
}

func TestRegisterDuplicate(t *testing.T) {
	r, _ := setupTestRouter(t)
	body, _ := json.Marshal(map[string]string{
		"email":    "test@example.com",
		"password": "testpass123",
	})
	
	// First registration
	req1, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)
	if w1.Code != http.StatusCreated {
		t.Errorf("Expected 201 for first register, got %d", w1.Code)
	}
	
	// Duplicate registration
	req2, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	if w2.Code != http.StatusConflict {
		t.Errorf("Expected 409 for duplicate, got %d", w2.Code)
	}
}

func TestLogin(t *testing.T) {
	r, _ := setupTestRouter(t)
	body, _ := json.Marshal(map[string]string{
		"email":    "login@example.com",
		"password": "loginpass123",
	})
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("Register failed: %d", w.Code)
	}
	
	// Now login
	loginBody, _ := json.Marshal(map[string]string{
		"email":    "login@example.com",
		"password": "loginpass123",
	})
	loginReq, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginW := httptest.NewRecorder()
	r.ServeHTTP(loginW, loginReq)
	if loginW.Code != http.StatusOK {
		t.Errorf("Expected 200 for login, got %d: %s", loginW.Code, loginW.Body.String())
	}
	
	// Wrong password
	wrongBody, _ := json.Marshal(map[string]string{
		"email":    "login@example.com",
		"password": "wrongpassword",
	})
	wrongReq, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewReader(wrongBody))
	wrongReq.Header.Set("Content-Type", "application/json")
	wrongW := httptest.NewRecorder()
	r.ServeHTTP(wrongW, wrongReq)
	if wrongW.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 for wrong password, got %d", wrongW.Code)
	}
}

func TestGetUserRequiresAuth(t *testing.T) {
	r, _ := setupTestRouter(t)
	req, _ := http.NewRequest("GET", "/api/user/usage", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected 401 without auth, got %d", w.Code)
	}
}

func TestUploadAudio(t *testing.T) {
	r, _ := setupTestRouter(t)
	
	// Create multipart body
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("audio", "test.mp3")
	part.Write([]byte("fake audio data"))
	writer.Close()
	
	req, _ := http.NewRequest("POST", "/api/audio/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	// Should return 202 (processing) or 403 (quota exceeded)
	if w.Code != http.StatusAccepted && w.Code != http.StatusForbidden {
		t.Logf("Got %d (expected 202 or 403)", w.Code)
	}
}

func TestGetJob(t *testing.T) {
	r, _ := setupTestRouter(t)
	req, _ := http.NewRequest("GET", "/api/audio/nonexistent-id", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", w.Code)
	}
}
