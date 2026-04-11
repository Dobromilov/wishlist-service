package domain

import "time"

type WishlistItem struct {
	ID          int       `json:"id" db:"id"`
	WishlistID  int       `json:"wishlist_id" db:"wishlist_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	URL         string    `json:"url,omitempty" db:"url"`
	Priority    *int      `json:"priority,omitempty" db:"priority"`
	ReservedBy  *int      `json:"reserved_by,omitempty" db:"reserved_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}
