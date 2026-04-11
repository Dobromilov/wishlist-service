package storage

import (
	"context"
	"errors"

	"wishlist-service/internal/domain"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository struct {
	db *Storage
}

func NewUserRepository(db *Storage) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.DB.QueryRow(ctx, query, user.Email, user.PasswordHash).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE email = $1`
	err := r.db.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, email, password_hash, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.DB.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
