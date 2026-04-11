package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"wishlist-service/internal/service"
	"wishlist-service/internal/validator"
)

type AuthHandler struct {
	authService *service.AuthService
	logger      *slog.Logger
}

func NewAuthHandler(authService *service.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req service.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Email(req.Email); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validator.Password(req.Password); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.authService.Register(r.Context(), req)
	if err != nil {
		if err == service.ErrEmailTaken {
			h.respondError(w, http.StatusConflict, "email already taken")
			return
		}
		h.logger.Error("register failed", "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to register")
		return
	}

	h.respondJSON(w, http.StatusCreated, service.AuthResponse{Token: token})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req service.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.authService.Login(r.Context(), req)
	if err != nil {
		h.logger.Info("login failed", "error", err)
		h.respondError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	h.respondJSON(w, http.StatusOK, service.AuthResponse{Token: token})
}

func (h *AuthHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *AuthHandler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
