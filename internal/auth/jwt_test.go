package auth

import (
	"os"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateAndValidateToken(t *testing.T) {
	tokenStr, err := GenerateToken("test-user-123")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	userID, err := ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}
	if userID != "test-user-123" {
		t.Errorf("Expected user_id 'test-user-123', got '%s'", userID)
	}
}

func TestExtractTokenFromHeader(t *testing.T) {
	tokenStr, _ := GenerateToken("user1")
	header := "Bearer " + tokenStr
	extracted, err := ExtractTokenFromHeader(header)
	if err != nil {
		t.Fatalf("Failed to extract token: %v", err)
	}
	if extracted != tokenStr {
		t.Errorf("Expected extracted token to match")
	}

	_, err = ExtractTokenFromHeader("InvalidFormat")
	if err == nil {
		t.Error("Expected error for invalid format")
	}
}

func TestValidateTamperedToken(t *testing.T) {
	secret := "test-secret"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": "test",
		"exp":     9999999999,
	})
	tokenStr, _ := token.SignedString([]byte(secret))

	os.Setenv("JWT_SECRET", secret)
	defer os.Unsetenv("JWT_SECRET")

	// Tamper with the token
	parts := split(tokenStr, ".")
	if len(parts) == 3 {
		parts[2] = "tampered"
		tamperedToken := parts[0] + "." + parts[1] + "." + parts[2]
		_, err := ValidateToken(tamperedToken)
		if err == nil {
			t.Error("Expected error for tampered token")
		}
	}
}
