package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"wishlist-service/internal/service"
)

type WishlistItemHandler struct {
	wishlistService *service.WishlistService
	logger          *slog.Logger
}

func NewWishlistItemHandler(wishlistService *service.WishlistService, logger *slog.Logger) *WishlistItemHandler {
	return &WishlistItemHandler{
		wishlistService: wishlistService,
		logger:          logger,
	}
}

func (h *WishlistItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	wishlistIDStr := strings.TrimPrefix(r.URL.Path, "/api/wishlists/")
	wishlistIDStr = strings.TrimSuffix(wishlistIDStr, "/items")
	wishlistID, err := strconv.Atoi(wishlistIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	var body struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Priority    *int   `json:"priority"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	item, err := h.wishlistService.AddItem(r.Context(), wishlistID, body.Name, body.Description, body.URL, body.Priority)
	if err != nil {
		h.logger.Error("create item failed", "error", err)
		h.respondError(w, http.StatusInternalServerError, "failed to create item")
		return
	}

	h.respondJSON(w, http.StatusCreated, item)
}

func (h *WishlistItemHandler) GetByWishlist(w http.ResponseWriter, r *http.Request) {
	wishlistIDStr := strings.TrimPrefix(r.URL.Path, "/api/wishlists/")
	wishlistIDStr = strings.TrimSuffix(wishlistIDStr, "/items")
	wishlistID, err := strconv.Atoi(wishlistIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	items, err := h.wishlistService.GetItems(r.Context(), wishlistID)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "failed to get items")
		return
	}

	h.respondJSON(w, http.StatusOK, items)
}

func (h *WishlistItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/wishlists/")
	parts := strings.Split(path, "/")
	if len(parts) != 4 {
		h.respondError(w, http.StatusBadRequest, "invalid url")
		return
	}

	wishlistID, err := strconv.Atoi(parts[0])
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	itemID, err := strconv.Atoi(parts[3])
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid item id")
		return
	}

	if err := h.wishlistService.DeleteItem(r.Context(), itemID, wishlistID); err != nil {
		h.respondError(w, http.StatusNotFound, "item not found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *WishlistItemHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *WishlistItemHandler) respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
