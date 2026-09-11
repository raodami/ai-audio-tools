package store

import (
	"testing"

	"github.com/google/uuid"
)

func newTestDB(t *testing.T) *Store {
	t.Helper()
	store, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test DB: %v", err)
	}
	return store
}

func TestCreateUser(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	id := uuid.New().String()
	err := s.CreateUser(id, "test@example.com", "hashed_password")
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	user, err := s.GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.Email)
	}
	if user.ID != id {
		t.Errorf("Expected ID %s, got %s", id, user.ID)
	}
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	id1 := uuid.New().String()
	id2 := uuid.New().String()

	err := s.CreateUser(id1, "dup@example.com", "pass1")
	if err != nil {
		t.Fatalf("First createUser failed: %v", err)
	}

	err = s.CreateUser(id2, "dup@example.com", "pass2")
	if err == nil {
		t.Fatal("Expected error for duplicate email, got nil")
	}
}

func TestIncrementUserUsage(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	id := uuid.New().String()
	s.CreateUser(id, "user@example.com", "pass")

	err := s.IncrementUserUsage(id, 10)
	if err != nil {
		t.Fatalf("IncrementUserUsage failed: %v", err)
	}

	user, _ := s.GetUserByID(id)
	if user.UsageMin != 10 {
		t.Errorf("Expected 10 minutes, got %d", user.UsageMin)
	}
}

func TestGetOrCreateSession(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	hash := "test_cookie_hash"

	sess1, err := s.GetOrCreateSession(hash, "1.2.3.4")
	if err != nil {
		t.Fatalf("GetOrCreateSession failed: %v", err)
	}
	if sess1.UsageMin != 0 {
		t.Errorf("Expected 0 minutes, got %d", sess1.UsageMin)
	}

	// Should return same session
	sess2, err := s.GetOrCreateSession(hash, "1.2.3.4")
	if err != nil {
		t.Fatalf("Second GetOrCreateSession failed: %v", err)
	}
	if sess1.ID != sess2.ID {
		t.Error("Expected same session ID")
	}
}

func TestIncrementSessionUsage(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	hash := "session_hash"
	s.GetOrCreateSession(hash, "1.2.3.4")

	err := s.IncrementSessionUsage(hash, 5)
	if err != nil {
		t.Fatalf("IncrementSessionUsage failed: %v", err)
	}

	sess, _ := s.GetOrCreateSession(hash, "1.2.3.4")
	if sess.UsageMin != 5 {
		t.Errorf("Expected 5 minutes, got %d", sess.UsageMin)
	}
}

func TestCheckQuota_FreeUser(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	id := uuid.New().String()
	s.CreateUser(id, "free@example.com", "pass")

	allowed, remaining, err := s.CheckQuota(id, "")
	if err != nil {
		t.Fatalf("CheckQuota failed: %v", err)
	}
	if !allowed {
		t.Error("Expected allowed=true for new user")
	}
	if remaining != FreeQuota {
		t.Errorf("Expected %d remaining, got %d", FreeQuota, remaining)
	}
}

func TestCheckQuota_Exceeded(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	id := uuid.New().String()
	s.CreateUser(id, "exceeded@example.com", "pass")
	s.IncrementUserUsage(id, FreeQuota+10)

	allowed, _, err := s.CheckQuota(id, "")
	if err != nil {
		t.Fatalf("CheckQuota failed: %v", err)
	}
	if allowed {
		t.Error("Expected allowed=false when quota exceeded")
	}
}

func TestCheckQuota_ProUser(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	id := uuid.New().String()
	s.CreateUser(id, "pro@example.com", "pass")

	// Mark as pro
	s.db.Exec("UPDATE users SET is_pro = 1 WHERE id = ?", id)

	allowed, remaining, err := s.CheckQuota(id, "")
	if err != nil {
		t.Fatalf("CheckQuota failed: %v", err)
	}
	if !allowed {
		t.Error("Expected allowed=true for pro user")
	}
	if remaining != -1 {
		t.Errorf("Expected -1 (unlimited) for pro user, got %d", remaining)
	}
}

func TestCheckQuota_Anonymous(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	hash := "anon_cookie"
	allowed, remaining, err := s.CheckQuota("", hash)
	if err != nil {
		t.Fatalf("CheckQuota failed: %v", err)
	}
	if !allowed {
		t.Error("Expected allowed=true for new anonymous user")
	}
	if remaining != FreeQuota {
		t.Errorf("Expected %d remaining, got %d", FreeQuota, remaining)
	}
}

func TestCreateJob(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	id := uuid.New().String()
	err := s.CreateJob(id, "user123", "", "test.mp3", 5000000)
	if err != nil {
		t.Fatalf("CreateJob failed: %v", err)
	}

	job, err := s.GetJob(id)
	if err != nil {
		t.Fatalf("GetJob failed: %v", err)
	}
	if job["status"] != "pending" {
		t.Errorf("Expected status 'pending', got %v", job["status"])
	}
	if job["file_name"] != "test.mp3" {
		t.Errorf("Expected file_name 'test.mp3', got %v", job["file_name"])
	}
}

func TestUpdateJobStatus(t *testing.T) {
	s := newTestDB(t)
	defer s.Close()

	id := uuid.New().String()
	s.CreateJob(id, "", "", "test.wav", 1000000)

	err := s.UpdateJobStatus(id, "completed", `{"text":"hello world"}`)
	if err != nil {
		t.Fatalf("UpdateJobStatus failed: %v", err)
	}

	job, _ := s.GetJob(id)
	if job["status"] != "completed" {
		t.Errorf("Expected status 'completed', got %v", job["status"])
	}
}

func TestFreeQuota_Constant(t *testing.T) {
	if FreeQuota != 30 {
		t.Errorf("Expected FreeQuota to be 30, got %d", FreeQuota)
	}
}
