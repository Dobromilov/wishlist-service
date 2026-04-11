package validator

import "testing"

func TestEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{"valid", "user@example.com", false},
		{"valid with dots", "first.last@example.com", false},
		{"invalid no @", "userexample.com", true},
		{"invalid empty", "", true},
		{"invalid no domain", "user@", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Email(tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("Email() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"valid", "password123", false},
		{"min length", "123456", false},
		{"too short", "12345", true},
		{"empty", "", true},
		{"spaces only", "   ", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Password(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Password() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTitle(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		wantErr bool
	}{
		{"valid", "Birthday", false},
		{"empty", "", true},
		{"spaces", "   ", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Title(tt.title)
			if (err != nil) != tt.wantErr {
				t.Errorf("Title() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPriority(t *testing.T) {
	ptr := func(v int) *int { return &v }
	tests := []struct {
		name     string
		priority *int
		wantErr  bool
	}{
		{"valid 1", ptr(1), false},
		{"valid 5", ptr(5), false},
		{"valid 10", ptr(10), false},
		{"nil", nil, false},
		{"too low", ptr(0), true},
		{"too high", ptr(11), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Priority(tt.priority)
			if (err != nil) != tt.wantErr {
				t.Errorf("Priority() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestItemName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "Headphones", false},
		{"empty", "", true},
		{"spaces", "   ", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ItemName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ItemName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
