package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"wishlist-service/internal/domain"
)

var ErrWishlistNotFound = errors.New("wishlist not found")

type WishlistRepository struct {
	db *pgxpool.Pool
}

func NewWishlistRepository(db *pgxpool.Pool) *WishlistRepository {
	return &WishlistRepository{db: db}
}

func (r *WishlistRepository) Create(ctx context.Context, wishlist *domain.Wishlist) error {
	query := `
		INSERT INTO wishlists (user_id, title, description, event_date, public_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query, wishlist.UserID, wishlist.Title, wishlist.Description, wishlist.EventDate, wishlist.PublicToken).
		Scan(&wishlist.ID, &wishlist.CreatedAt, &wishlist.UpdatedAt)
	return err
}

func (r *WishlistRepository) GetByID(ctx context.Context, id int) (*domain.Wishlist, error) {
	w := &domain.Wishlist{}
	query := `SELECT id, user_id, title, description, event_date, public_token, created_at, updated_at FROM wishlists WHERE id = $1`
	err := r.db.QueryRow(ctx, query, id).
		Scan(&w.ID, &w.UserID, &w.Title, &w.Description, &w.EventDate, &w.PublicToken, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, ErrWishlistNotFound
	}
	return w, nil
}

func (r *WishlistRepository) GetByPublicToken(ctx context.Context, token string) (*domain.Wishlist, error) {
	w := &domain.Wishlist{}
	query := `SELECT id, user_id, title, description, event_date, public_token, created_at, updated_at FROM wishlists WHERE public_token = $1`
	err := r.db.QueryRow(ctx, query, token).
		Scan(&w.ID, &w.UserID, &w.Title, &w.Description, &w.EventDate, &w.PublicToken, &w.CreatedAt, &w.UpdatedAt)
	if err != nil {
		return nil, ErrWishlistNotFound
	}
	return w, nil
}

func (r *WishlistRepository) GetByUserID(ctx context.Context, userID int) ([]domain.Wishlist, error) {
	query := `SELECT id, user_id, title, description, event_date, public_token, created_at, updated_at FROM wishlists WHERE user_id = $1`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wishlists []domain.Wishlist
	for rows.Next() {
		var w domain.Wishlist
		if err := rows.Scan(&w.ID, &w.UserID, &w.Title, &w.Description, &w.EventDate, &w.PublicToken, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		wishlists = append(wishlists, w)
	}
	return wishlists, nil
}

func (r *WishlistRepository) Update(ctx context.Context, wishlist *domain.Wishlist) error {
	query := `UPDATE wishlists SET title = $1, description = $2, event_date = $3, updated_at = NOW() WHERE id = $4 AND user_id = $5`
	res, err := r.db.Exec(ctx, query, wishlist.Title, wishlist.Description, wishlist.EventDate, wishlist.ID, wishlist.UserID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrWishlistNotFound
	}
	return nil
}

func (r *WishlistRepository) Delete(ctx context.Context, id, userID int) error {
	query := `DELETE FROM wishlists WHERE id = $1 AND user_id = $2`
	res, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return ErrWishlistNotFound
	}
	return nil
}
