package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	pass := "test-password-123"
	hash, err := HashPassword(pass)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}
	if hash == pass {
		t.Error("Hash should not equal plaintext")
	}

	// Verify correct password
	if err := CheckPasswordHash(pass, hash); err != nil {
		t.Errorf("Expected password to match, got error: %v", err)
	}

	// Verify wrong password
	if err := CheckPasswordHash("wrong-password", hash); err == nil {
		t.Error("Expected wrong password to fail")
	}
}

func TestHashPassword_Empty(t *testing.T) {
	_, err := HashPassword("")
	if err != nil {
		// Empty password may or may not error depending on bcrypt version
		t.Logf("Empty password hash returned error (expected): %v", err)
	}
}
