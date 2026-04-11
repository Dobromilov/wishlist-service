package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"wishlist-service/internal/domain"
)

var ErrItemNotFound = errors.New("wishlist item not found")

type WishlistItemRepository struct {
	db *pgxpool.Pool
}

func NewWishlistItemRepository(db *pgxpool.Pool) *WishlistItemRepository {
	return &WishlistItemRepository{db: db}
}

func (r *WishlistItemRepository) Create(ctx context.Context, item *domain.WishlistItem) error {
	query := `
		INSERT INTO wishlist_items (wishlist_id, name, description, url, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, item.WishlistID, item.Name, item.Description, item.URL, item.Priority).
		Scan(&item.ID, &item.CreatedAt, &item.UpdatedAt)
	return err
}

func (r *WishlistItemRepository) GetByID(ctx context.Context, id int) (*domain.WishlistItem, error) {
	item := &domain.WishlistItem{}
	query := `SELECT id, wishlist_id, name, description, url, priority, reserved_by, created_at, updated_at FROM wishlist_items WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).
		Scan(&item.ID, &item.WishlistID, &item.Name, &item.Description, &item.URL, &item.Priority, &item.ReservedBy, &item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, ErrItemNotFound
	}
	return item, nil
}

func (r *WishlistItemRepository) GetByWishlistID(ctx context.Context, wishlistID int) ([]domain.WishlistItem, error) {
	query := `SELECT id, wishlist_id, name, description, url, priority, reserved_by, created_at, updated_at FROM wishlist_items WHERE wishlist_id = $1`
	rows, err := r.db.Query(ctx, query, wishlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.WishlistItem
	for rows.Next() {
		var item domain.WishlistItem
		if err := rows.Scan(&item.ID, &item.WishlistID, &item.Name, &item.Description, &item.URL, &item.Priority, &item.ReservedBy, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *WishlistItemRepository) Update(ctx context.Context, item *domain.WishlistItem) error {
	query := `UPDATE wishlist_items SET name = $1, description = $2, url = $3, priority = $4, updated_at = NOW() WHERE id = $5 AND wishlist_id = $6`
	res, err := r.db.Exec(ctx, query, item.Name, item.Description, item.URL, item.Priority, item.ID, item.WishlistID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrItemNotFound
	}
	return nil
}

func (r *WishlistItemRepository) Delete(ctx context.Context, id, wishlistID int) error {
	query := `DELETE FROM wishlist_items WHERE id = $1 AND wishlist_id = $2`
	res, err := r.db.Exec(ctx, query, id, wishlistID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrItemNotFound
	}
	return nil
}

func (r *WishlistItemRepository) Reserve(ctx context.Context, itemID int) error {
	query := `UPDATE wishlist_items SET reserved_by = 1, updated_at = NOW() WHERE id = $1 AND reserved_by IS NULL`
	res, err := r.db.Exec(ctx, query, itemID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("item already reserved")
	}
	return nil
}
