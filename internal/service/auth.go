package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"wishlist-service/internal/domain"
	"wishlist-service/internal/storage"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrEmailTaken         = errors.New("email already taken")
)

type AuthService struct {
	userRepo  *storage.UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo *storage.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &domain.User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			return "", ErrEmailTaken
		}
		return "", err
	}

	return s.generateToken(user.ID, user.Email)
}

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	return s.generateToken(user.ID, user.Email)
}

func (s *AuthService) ValidateToken(tokenStr string) (int, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return s.jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidToken
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, ErrInvalidToken
	}

	return int(userID), nil
}

func (s *AuthService) generateToken(userID int, email string) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(2 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
