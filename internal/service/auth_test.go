package service

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func TestGenerateAndValidateToken(t *testing.T) {
	secret := "test-secret-key"
	svc := &AuthService{jwtSecret: []byte(secret)}

	token, err := svc.generateToken(42, "user@example.com")
	if err != nil {
		t.Fatalf("generateToken() error = %v", err)
	}
	if token == "" {
		t.Fatal("generateToken() returned empty token")
	}

	userID, err := svc.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken() error = %v", err)
	}
	if userID != 42 {
		t.Errorf("ValidateToken() userID = %d, want 42", userID)
	}
}

func TestValidateToken_Invalid(t *testing.T) {
	svc := &AuthService{jwtSecret: []byte("secret")}

	_, err := svc.ValidateToken("totally-invalid-token")
	if err == nil {
		t.Fatal("ValidateToken() expected error for invalid token")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	svc1 := &AuthService{jwtSecret: []byte("secret1")}
	svc2 := &AuthService{jwtSecret: []byte("secret2")}

	token, err := svc1.generateToken(1, "test@test.com")
	if err != nil {
		t.Fatalf("generateToken() error = %v", err)
	}

	_, err = svc2.ValidateToken(token)
	if err == nil {
		t.Fatal("ValidateToken() expected error for wrong secret")
	}
}

func TestValidateToken_Expired(t *testing.T) {
	svc := &AuthService{jwtSecret: []byte("secret")}

	claims := jwt.MapClaims{
		"user_id": 1,
		"email":   "test@test.com",
		"exp":     1000000000,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, _ := token.SignedString(svc.jwtSecret)

	_, err := svc.ValidateToken(expiredToken)
	if err == nil {
		t.Fatal("ValidateToken() expected error for expired token")
	}
}

func TestBcryptHashAndCompare(t *testing.T) {
	password := "mysecurepassword123"

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	err = bcrypt.CompareHashAndPassword(hashed, []byte(password))
	if err != nil {
		t.Errorf("CompareHashAndPassword() error = %v", err)
	}

	err = bcrypt.CompareHashAndPassword(hashed, []byte("wrongpassword"))
	if err == nil {
		t.Fatal("CompareHashAndPassword() expected error for wrong password")
	}
}
