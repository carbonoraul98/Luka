package entity_test

import (
	"errors"
	"testing"

	"github.com/raulcarbono/luka/backend/internal/domain/entity"
)

func TestNewUser(t *testing.T) {
	validID := "user-123"
	validEmail := "test@example.com"
	validPasswordHash := "$2a$10$hashedpasswordstring"
	validCurrency := "USD"

	tests := []struct {
		name         string
		id           string
		email        string
		passwordHash string
		currency     string
		wantErr      error
	}{
		{
			name:         "valid user with USD",
			id:           validID,
			email:        validEmail,
			passwordHash: validPasswordHash,
			currency:     "USD",
			wantErr:      nil,
		},
		{
			name:         "valid user with COP lowercase trimmed",
			id:           validID,
			email:        "  user@domain.co  ",
			passwordHash: validPasswordHash,
			currency:     " cop ",
			wantErr:      nil,
		},
		{
			name:         "empty id",
			id:           "   ",
			email:        validEmail,
			passwordHash: validPasswordHash,
			currency:     validCurrency,
			wantErr:      entity.ErrEmptyUserID,
		},
		{
			name:         "invalid email without at",
			id:           validID,
			email:        "invalid-email",
			passwordHash: validPasswordHash,
			currency:     validCurrency,
			wantErr:      entity.ErrInvalidEmail,
		},
		{
			name:         "invalid email format with name display",
			id:           validID,
			email:        "John Doe <john@example.com>",
			passwordHash: validPasswordHash,
			currency:     validCurrency,
			wantErr:      entity.ErrInvalidEmail,
		},
		{
			name:         "empty password hash",
			id:           validID,
			email:        validEmail,
			passwordHash: "   ",
			currency:     validCurrency,
			wantErr:      entity.ErrEmptyPasswordHash,
		},
		{
			name:         "unsupported currency",
			id:           validID,
			email:        validEmail,
			passwordHash: validPasswordHash,
			currency:     "XYZ",
			wantErr:      entity.ErrUnsupportedCurrency,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user, err := entity.NewUser(tt.id, tt.email, tt.passwordHash, tt.currency)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				if user != nil {
					t.Fatalf("expected nil user on error, got %+v", user)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if user.ID != tt.id {
				t.Errorf("got ID %q, want %q", user.ID, tt.id)
			}
			if user.PasswordHash != tt.passwordHash {
				t.Errorf("got PasswordHash %q, want %q", user.PasswordHash, tt.passwordHash)
			}
			if user.IsDeleted() {
				t.Error("new user should not be deleted")
			}
		})
	}
}

func TestUser_UpdatePreferences(t *testing.T) {
	user, err := entity.NewUser("id-1", "user@test.com", "hash", "USD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("valid currency update", func(t *testing.T) {
		err := user.UpdatePreferences(entity.UserPreferences{BaseCurrency: "COP"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if user.Preferences.BaseCurrency != "COP" {
			t.Errorf("got BaseCurrency %q, want COP", user.Preferences.BaseCurrency)
		}
	})

	t.Run("invalid currency update", func(t *testing.T) {
		err := user.UpdatePreferences(entity.UserPreferences{BaseCurrency: "ARS"})
		if !errors.Is(err, entity.ErrUnsupportedCurrency) {
			t.Fatalf("expected ErrUnsupportedCurrency, got %v", err)
		}
	})

	t.Run("cannot update if deleted", func(t *testing.T) {
		if err := user.Delete(); err != nil {
			t.Fatalf("failed to delete user: %v", err)
		}
		err := user.UpdatePreferences(entity.UserPreferences{BaseCurrency: "EUR"})
		if !errors.Is(err, entity.ErrUserAlreadyDeleted) {
			t.Fatalf("expected ErrUserAlreadyDeleted, got %v", err)
		}
	})
}

