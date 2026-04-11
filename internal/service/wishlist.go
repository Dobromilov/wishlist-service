package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"wishlist-service/internal/domain"
)

var ErrForbidden = errors.New("forbidden")

type WishlistRepository interface {
	Create(ctx context.Context, wishlist *domain.Wishlist) error
	GetByID(ctx context.Context, id int) (*domain.Wishlist, error)
	GetByPublicToken(ctx context.Context, token string) (*domain.Wishlist, error)
	GetByUserID(ctx context.Context, userID int) ([]domain.Wishlist, error)
	Update(ctx context.Context, wishlist *domain.Wishlist) error
	Delete(ctx context.Context, id, userID int) error
}

type WishlistItemRepository interface {
	Create(ctx context.Context, item *domain.WishlistItem) error
	GetByID(ctx context.Context, id int) (*domain.WishlistItem, error)
	GetByWishlistID(ctx context.Context, wishlistID int) ([]domain.WishlistItem, error)
	Update(ctx context.Context, item *domain.WishlistItem) error
	Delete(ctx context.Context, id, wishlistID int) error
	Reserve(ctx context.Context, itemID int) error
}

type WishlistService struct {
	wishlistRepo WishlistRepository
	itemRepo     WishlistItemRepository
}

func NewWishlistService(wishlistRepo WishlistRepository, itemRepo WishlistItemRepository) *WishlistService {
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
	if eventDate != nil && *eventDate != "" {
		t, err := time.Parse("2006-01-02", *eventDate)
		if err == nil {
			w.EventDate = &t
		}
	}
	if err := s.wishlistRepo.Create(ctx, w); err != nil {
		return nil, err
	}
	return w, nil
}

func (s *WishlistService) GetByID(ctx context.Context, id, userID int) (*domain.Wishlist, error) {
	w, err := s.wishlistRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if w.UserID != userID {
		return nil, ErrForbidden
	}
	return w, nil
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

func (s *WishlistService) OwnsWishlist(ctx context.Context, wishlistID, userID int) error {
	w, err := s.wishlistRepo.GetByID(ctx, wishlistID)
	if err != nil {
		return err
	}
	if w.UserID != userID {
		return ErrForbidden
	}
	return nil
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
