package domain

import (
	"time"
)

type Wishlist struct {
	ID          int        `json:"id" db:"id"`
	UserID      int        `json:"user_id" db:"user_id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description,omitempty" db:"description"`
	EventDate   *time.Time `json:"event_date,omitempty" db:"event_date"`
	PublicToken string     `json:"public_token" db:"public_token"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}
