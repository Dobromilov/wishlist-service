package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"wishlist-service/internal/domain"
	"wishlist-service/internal/storage"
)

type WishlistService struct {
	wishlistRepo *storage.WishlistRepository
	itemRepo     *storage.WishlistItemRepository
}

func NewWishlistService(wishlistRepo *storage.WishlistRepository, itemRepo *storage.WishlistItemRepository) *WishlistService {
	return &WishlistService{
		wishlistRepo: wishlistRepo,
		itemRepo:     itemRepo,
	}
}

func (s *WishlistService) Create(ctx context.Context, userID int, title string, description string, eventDate *string) (*domain.Wishlist, error) {
	token := generateToken()
	w := &domain.Wishlist{
		UserID:      userID,
		Title:       title,
		Description: description,
		PublicToken: token,
	}
	if err := s.wishlistRepo.Create(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WishlistService) GetByID(ctx context.Context, id, userID int) (*domain.Wishlist, error) {
	return s.wishlistRepo.GetByID(ctx, id)
}

func (s *WishlistService) GetByUser(ctx context.Context, userID int) ([]domain.Wishlist, error) {
	return s.wishlistRepo.GetByUserID(ctx, userID)
}

func (s *WishlistService) Update(ctx context.Context, userID int, w *domain.Wishlist) error {
	return s.wishlistRepo.Update(ctx, w)
}

func (s *WishlistService) Delete(ctx context.Context, id, userID int) error {
	return s.wishlistRepo.Delete(ctx, id, userID)
}

func (s *WishlistService) GetByPublicToken(ctx context.Context, token string) (*domain.Wishlist, error) {
	return s.wishlistRepo.GetByPublicToken(ctx, token)
}

func (s *WishlistService) AddItem(ctx context.Context, wishlistID int, name, description, url string, priority *int) (*domain.WishlistItem, error) {
	item := &domain.WishlistItem{
		WishlistID:  wishlistID,
		Name:        name,
		Description: description,
		URL:         url,
		Priority:    priority,
	}
	if err := s.itemRepo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *WishlistService) GetItems(ctx context.Context, wishlistID int) ([]domain.WishlistItem, error) {
	return s.itemRepo.GetByWishlistID(ctx, wishlistID)
}

func (s *WishlistService) UpdateItem(ctx context.Context, item *domain.WishlistItem) error {
	return s.itemRepo.Update(ctx, item)
}

func (s *WishlistService) DeleteItem(ctx context.Context, id, wishlistID int) error {
	return s.itemRepo.Delete(ctx, id, wishlistID)
}

func (s *WishlistService) ReserveItem(ctx context.Context, itemID int) error {
	return s.itemRepo.Reserve(ctx, itemID)
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
