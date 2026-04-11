package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"wishlist-service/internal/service"
	"wishlist-service/internal/validator"
)

type WishlistHandler struct {
	wishlistService *service.WishlistService
	logger          *slog.Logger
}

func NewWishlistHandler(wishlistService *service.WishlistService, logger *slog.Logger) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
		logger:          logger,
	}
}

func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var body struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		EventDate   *string `json:"event_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := validator.Title(body.Title); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	wishlist, err := h.wishlistService.Create(r.Context(), userID, body.Title, body.Description, body.EventDate)
	if err != nil {
		h.logger.Error("create wishlist failed", "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to create wishlist")
		return
	}

	h.respondJSON(w, http.StatusCreated, wishlist)
}

func (h *WishlistHandler) GetMy(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlists, err := h.wishlistService.GetByUser(r.Context(), userID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "failed to get wishlists")
		return
	}

	h.respondJSON(w, http.StatusOK, wishlists)
}

func (h *WishlistHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	var body struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		EventDate   *string `json:"event_date"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	wishlist, err := h.wishlistService.GetByID(r.Context(), id, userID)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "wishlist not found")
		return
	}

	wishlist.Title = body.Title
	wishlist.Description = body.Description
	if body.EventDate != nil && *body.EventDate != "" {
		t, err := time.Parse("2006-01-02", *body.EventDate)
		if err == nil {
			wishlist.EventDate = &t
		}
	}

	if err := h.wishlistService.Update(r.Context(), userID, wishlist); err != nil {
		h.respondError(w, http.StatusInternalServerError, "failed to update wishlist")
		return
	}

	h.respondJSON(w, http.StatusOK, wishlist)
}

func (h *WishlistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	if err := h.wishlistService.Delete(r.Context(), id, userID); err != nil {
		h.respondError(w, http.StatusNotFound, "wishlist not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WishlistHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *WishlistHandler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
