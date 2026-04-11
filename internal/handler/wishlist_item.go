package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"wishlist-service/internal/service"
	"wishlist-service/internal/validator"
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
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlistID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	if err := h.wishlistService.OwnsWishlist(r.Context(), wishlistID, userID); err != nil {
		h.respondError(w, http.StatusForbidden, "forbidden")
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

	if err := validator.ItemName(body.Name); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := validator.Priority(body.Priority); err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
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
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlistID, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	if err := h.wishlistService.OwnsWishlist(r.Context(), wishlistID, userID); err != nil {
		h.respondError(w, http.StatusForbidden, "forbidden")
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
	userID, ok := GetUserIDFromContext(r.Context())
	if !ok {
		h.respondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	wishlistID, err := strconv.Atoi(r.PathValue("wishlistId"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid wishlist id")
		return
	}

	if err := h.wishlistService.OwnsWishlist(r.Context(), wishlistID, userID); err != nil {
		h.respondError(w, http.StatusForbidden, "forbidden")
		return
	}

	itemID, err := strconv.Atoi(r.PathValue("itemId"))
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
