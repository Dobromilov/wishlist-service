package service

import (
	"context"
	"errors"
	"testing"

	"wishlist-service/internal/domain"
)

var (
	errWishlistNotFound = errors.New("wishlist not found")
	errItemNotFound     = errors.New("wishlist item not found")
)

type testWishlistRepo struct {
	wishlists map[int]*domain.Wishlist
	items     map[int]*domain.WishlistItem
	nextWID   int
	nextIID   int
}

func newTestWishlistRepo() *testWishlistRepo {
	return &testWishlistRepo{
		wishlists: make(map[int]*domain.Wishlist),
		items:     make(map[int]*domain.WishlistItem),
		nextWID:   1,
		nextIID:   1,
	}
}

func (r *testWishlistRepo) Create(ctx context.Context, w *domain.Wishlist) error {
	w.ID = r.nextWID
	r.nextWID++
	r.wishlists[w.ID] = w
	return nil
}

func (r *testWishlistRepo) GetByID(ctx context.Context, id int) (*domain.Wishlist, error) {
	w, ok := r.wishlists[id]
	if !ok {
		return nil, errWishlistNotFound
	}
	return w, nil
}

func (r *testWishlistRepo) GetByPublicToken(ctx context.Context, token string) (*domain.Wishlist, error) {
	for _, w := range r.wishlists {
		if w.PublicToken == token {
			return w, nil
		}
	}
	return nil, errWishlistNotFound
}

func (r *testWishlistRepo) GetByUserID(ctx context.Context, userID int) ([]domain.Wishlist, error) {
	var result []domain.Wishlist
	for _, w := range r.wishlists {
		if w.UserID == userID {
			result = append(result, *w)
		}
	}
	return result, nil
}

func (r *testWishlistRepo) Update(ctx context.Context, w *domain.Wishlist) error {
	if _, ok := r.wishlists[w.ID]; !ok {
		return errWishlistNotFound
	}
	r.wishlists[w.ID] = w
	return nil
}

func (r *testWishlistRepo) Delete(ctx context.Context, id, userID int) error {
	w, ok := r.wishlists[id]
	if !ok || w.UserID != userID {
		return errWishlistNotFound
	}
	delete(r.wishlists, id)
	return nil
}

func (r *testWishlistRepo) CreateItem(ctx context.Context, item *domain.WishlistItem) error {
	item.ID = r.nextIID
	r.nextIID++
	r.items[item.ID] = item
	return nil
}

func (r *testWishlistRepo) GetItemByID(ctx context.Context, id int) (*domain.WishlistItem, error) {
	item, ok := r.items[id]
	if !ok {
		return nil, errItemNotFound
	}
	return item, nil
}

func (r *testWishlistRepo) GetItemsByWishlistID(ctx context.Context, wishlistID int) ([]domain.WishlistItem, error) {
	var result []domain.WishlistItem
	for _, item := range r.items {
		if item.WishlistID == wishlistID {
			result = append(result, *item)
		}
	}
	return result, nil
}

func (r *testWishlistRepo) UpdateItem(ctx context.Context, item *domain.WishlistItem) error {
	if _, ok := r.items[item.ID]; !ok {
		return errItemNotFound
	}
	r.items[item.ID] = item
	return nil
}

func (r *testWishlistRepo) DeleteItem(ctx context.Context, id, wishlistID int) error {
	item, ok := r.items[id]
	if !ok || item.WishlistID != wishlistID {
		return errItemNotFound
	}
	delete(r.items, id)
	return nil
}

func (r *testWishlistRepo) ReserveItem(ctx context.Context, itemID int) error {
	item, ok := r.items[itemID]
	if !ok {
		return errItemNotFound
	}
	if item.ReservedBy != nil {
		return errors.New("item already reserved")
	}
	one := 1
	item.ReservedBy = &one
	return nil
}

func TestWishlistService_Create(t *testing.T) {
	repo := newTestWishlistRepo()
	svc := NewWishlistService(repoWrapper{repo}, itemWrapper{repo})

	eventDate := "2025-12-25"
	w, err := svc.Create(context.Background(), 1, "Birthday", "My gifts", &eventDate)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if w.Title != "Birthday" {
		t.Errorf("Create() title = %s, want Birthday", w.Title)
	}
	if w.UserID != 1 {
		t.Errorf("Create() userID = %d, want 1", w.UserID)
	}
	if w.EventDate == nil {
		t.Error("Create() eventDate is nil")
	}
}

func TestWishlistService_GetByID_Ownership(t *testing.T) {
	repo := newTestWishlistRepo()
	svc := NewWishlistService(repoWrapper{repo}, itemWrapper{repo})

	svc.Create(context.Background(), 1, "My Wishlist", "", nil)

	_, err := svc.GetByID(context.Background(), 1, 2)
	if err != ErrForbidden {
		t.Errorf("GetByID() error = %v, want ErrForbidden", err)
	}
}

func TestWishlistService_OwnsWishlist(t *testing.T) {
	repo := newTestWishlistRepo()
	svc := NewWishlistService(repoWrapper{repo}, itemWrapper{repo})

	svc.Create(context.Background(), 1, "My Wishlist", "", nil)

	err := svc.OwnsWishlist(context.Background(), 1, 1)
	if err != nil {
		t.Errorf("OwnsWishlist() error = %v", err)
	}

	err = svc.OwnsWishlist(context.Background(), 1, 999)
	if err != ErrForbidden {
		t.Errorf("OwnsWishlist() error = %v, want ErrForbidden", err)
	}
}

func TestWishlistService_AddItem(t *testing.T) {
	repo := newTestWishlistRepo()
	svc := NewWishlistService(repoWrapper{repo}, itemWrapper{repo})

	svc.Create(context.Background(), 1, "My Wishlist", "", nil)

	priority := 5
	item, err := svc.AddItem(context.Background(), 1, "Headphones", "Nice headphones", "https://example.com", &priority)
	if err != nil {
		t.Fatalf("AddItem() error = %v", err)
	}
	if item.Name != "Headphones" {
		t.Errorf("AddItem() name = %s, want Headphones", item.Name)
	}
	if item.Priority == nil || *item.Priority != 5 {
		t.Errorf("AddItem() priority = %v, want 5", item.Priority)
	}
}

func TestWishlistService_GetByUser(t *testing.T) {
	repo := newTestWishlistRepo()
	svc := NewWishlistService(repoWrapper{repo}, itemWrapper{repo})

	svc.Create(context.Background(), 1, "List 1", "", nil)
	svc.Create(context.Background(), 1, "List 2", "", nil)
	svc.Create(context.Background(), 2, "Other List", "", nil)

	lists, err := svc.GetByUser(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetByUser() error = %v", err)
	}
	if len(lists) != 2 {
		t.Errorf("GetByUser() count = %d, want 2", len(lists))
	}
}

type repoWrapper struct {
	r *testWishlistRepo
}

func (w repoWrapper) Create(ctx context.Context, wishlist *domain.Wishlist) error {
	return w.r.Create(ctx, wishlist)
}

func (w repoWrapper) GetByID(ctx context.Context, id int) (*domain.Wishlist, error) {
	return w.r.GetByID(ctx, id)
}

func (w repoWrapper) GetByPublicToken(ctx context.Context, token string) (*domain.Wishlist, error) {
	return w.r.GetByPublicToken(ctx, token)
}

func (w repoWrapper) GetByUserID(ctx context.Context, userID int) ([]domain.Wishlist, error) {
	return w.r.GetByUserID(ctx, userID)
}

func (w repoWrapper) Update(ctx context.Context, wishlist *domain.Wishlist) error {
	return w.r.Update(ctx, wishlist)
}

func (w repoWrapper) Delete(ctx context.Context, id, userID int) error {
	return w.r.Delete(ctx, id, userID)
}

type itemWrapper struct {
	r *testWishlistRepo
}

func (w itemWrapper) Create(ctx context.Context, item *domain.WishlistItem) error {
	return w.r.CreateItem(ctx, item)
}

func (w itemWrapper) GetByID(ctx context.Context, id int) (*domain.WishlistItem, error) {
	return w.r.GetItemByID(ctx, id)
}

func (w itemWrapper) GetByWishlistID(ctx context.Context, wishlistID int) ([]domain.WishlistItem, error) {
	return w.r.GetItemsByWishlistID(ctx, wishlistID)
}

func (w itemWrapper) Update(ctx context.Context, item *domain.WishlistItem) error {
	return w.r.UpdateItem(ctx, item)
}

func (w itemWrapper) Delete(ctx context.Context, id, wishlistID int) error {
	return w.r.DeleteItem(ctx, id, wishlistID)
}

func (w itemWrapper) Reserve(ctx context.Context, itemID int) error {
	return w.r.ReserveItem(ctx, itemID)
}
