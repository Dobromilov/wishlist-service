package validator

import (
	"errors"
	"net/mail"
	"strings"
)

var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrShortPassword   = errors.New("password must be at least 6 characters")
	ErrEmptyTitle      = errors.New("title is required")
	ErrInvalidPriority = errors.New("priority must be between 1 and 10")
	ErrEmptyName       = errors.New("item name is required")
)

func Email(email string) error {
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
	}
	return nil
}

func Password(password string) error {
	if len(strings.TrimSpace(password)) < 6 {
		return ErrShortPassword
	}
	return nil
}

func Title(title string) error {
	if strings.TrimSpace(title) == "" {
		return ErrEmptyTitle
	}
	return nil
}

func Priority(p *int) error {
	if p == nil {
		return nil
	}
	if *p < 1 || *p > 10 {
		return ErrInvalidPriority
	}
	return nil
}

func ItemName(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrEmptyName
	}
	return nil
}
