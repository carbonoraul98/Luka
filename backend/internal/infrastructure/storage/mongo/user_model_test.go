package mongo_test

import (
	"testing"
	"time"

	"github.com/raulcarbono/luka/backend/internal/domain/entity"
	mongostorage "github.com/raulcarbono/luka/backend/internal/infrastructure/storage/mongo"
)

func TestUserModel_Mapping(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)
	domainUser := &entity.User{
		ID:           "usr-uuid-123",
		Email:        "architect@luka.app",
		PasswordHash: "$2a$10$hashedstring",
		Preferences: entity.UserPreferences{
			BaseCurrency: "USD",
		},
		CreatedAt: now,
		UpdatedAt: now,
		DeletedAt: nil,
	}

	doc := mongostorage.FromDomain(domainUser)

	if doc.ID != domainUser.ID {
		t.Errorf("got doc ID %q, want %q", doc.ID, domainUser.ID)
	}
	if doc.Email != domainUser.Email {
		t.Errorf("got doc Email %q, want %q", doc.Email, domainUser.Email)
	}
	if doc.PasswordHash != domainUser.PasswordHash {
		t.Errorf("got doc PasswordHash %q, want %q", doc.PasswordHash, domainUser.PasswordHash)
	}
	if doc.Preferences.BaseCurrency != domainUser.Preferences.BaseCurrency {
		t.Errorf("got doc BaseCurrency %q, want %q", doc.Preferences.BaseCurrency, domainUser.Preferences.BaseCurrency)
	}

	restoredUser := doc.ToDomain()

	if restoredUser.ID != domainUser.ID {
		t.Errorf("got restored ID %q, want %q", restoredUser.ID, domainUser.ID)
	}
	if restoredUser.Email != domainUser.Email {
		t.Errorf("got restored Email %q, want %q", restoredUser.Email, domainUser.Email)
	}
	if restoredUser.PasswordHash != domainUser.PasswordHash {
		t.Errorf("got restored PasswordHash %q, want %q", restoredUser.PasswordHash, domainUser.PasswordHash)
	}
	if restoredUser.Preferences.BaseCurrency != domainUser.Preferences.BaseCurrency {
		t.Errorf("got restored BaseCurrency %q, want %q", restoredUser.Preferences.BaseCurrency, domainUser.Preferences.BaseCurrency)
	}
	if !restoredUser.CreatedAt.Equal(domainUser.CreatedAt) {
		t.Errorf("got restored CreatedAt %v, want %v", restoredUser.CreatedAt, domainUser.CreatedAt)
	}
	if restoredUser.IsDeleted() {
		t.Error("restored user should not be deleted")
	}
}

