package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"wishlist-service/internal/service"
)

type PublicHandler struct {
	wishlistService *service.WishlistService
	logger          *slog.Logger
}

func NewPublicHandler(wishlistService *service.WishlistService, logger *slog.Logger) *PublicHandler {
	return &PublicHandler{
		wishlistService: wishlistService,
		logger:          logger,
	}
}

func (h *PublicHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.URL.Path, "/public/")

	wishlist, err := h.wishlistService.GetByPublicToken(r.Context(), token)
	if err != nil {
		h.respondError(w, http.StatusNotFound, "wishlist not found")
		return
	}

	items, err := h.wishlistService.GetItems(r.Context(), wishlist.ID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "failed to get items")
		return
	}

	response := map[string]interface{}{
		"wishlist": wishlist,
		"items":    items,
	}

	h.respondJSON(w, http.StatusOK, response)
}

func (h *PublicHandler) ReserveItem(w http.ResponseWriter, r *http.Request) {
	var body struct {
		ItemID int `json:"item_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.wishlistService.ReserveItem(r.Context(), body.ItemID); err != nil {
		h.respondError(w, http.StatusConflict, "item already reserved")
		return
	}

	h.respondJSON(w, http.StatusOK, map[string]string{"status": "reserved"})
}

func (h *PublicHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *PublicHandler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
