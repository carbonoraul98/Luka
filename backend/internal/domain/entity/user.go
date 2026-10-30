package entity

import (
	"errors"
	"net/mail"
	"strings"
	"time"
)

var (
	ErrEmptyUserID         = errors.New("user ID cannot be empty")
	ErrInvalidEmail        = errors.New("invalid email address")
	ErrEmptyPasswordHash   = errors.New("password hash cannot be empty")
	ErrUnsupportedCurrency = errors.New("unsupported base currency")
	ErrUserAlreadyDeleted  = errors.New("user is already deleted")
)

// SupportedCurrencies lists the accepted standard ISO currency codes.
var SupportedCurrencies = map[string]struct{}{
	"USD": {},
	"COP": {},
	"EUR": {},
}

// UserPreferences contains user-specific application preferences.
type UserPreferences struct {
	BaseCurrency string `json:"base_currency"`
}

// Validate checks whether preference attributes satisfy domain invariants.
func (p UserPreferences) Validate() error {
	currency := strings.ToUpper(strings.TrimSpace(p.BaseCurrency))
	if _, ok := SupportedCurrencies[currency]; !ok {
		return ErrUnsupportedCurrency
	}
	return nil
}

// User represents the core user domain entity.
type User struct {
	ID           string          `json:"id"`
	Email        string          `json:"email"`
	PasswordHash string          `json:"-"`
	Preferences  UserPreferences `json:"preferences"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    *time.Time      `json:"deleted_at,omitempty"`
}

// NewUser creates and validates a new User instance.
func NewUser(id, email, passwordHash, baseCurrency string) (*User, error) {
	cleanID := strings.TrimSpace(id)
	if cleanID == "" {
		return nil, ErrEmptyUserID
	}

	cleanEmail := strings.TrimSpace(email)
	parsedEmail, err := mail.ParseAddress(cleanEmail)
	if err != nil || parsedEmail.Address != cleanEmail {
		return nil, ErrInvalidEmail
	}

	cleanPasswordHash := strings.TrimSpace(passwordHash)
	if cleanPasswordHash == "" {
		return nil, ErrEmptyPasswordHash
	}

	currency := strings.ToUpper(strings.TrimSpace(baseCurrency))
	prefs := UserPreferences{BaseCurrency: currency}
	if err := prefs.Validate(); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &User{
		ID:           cleanID,
		Email:        cleanEmail,
		PasswordHash: cleanPasswordHash,
		Preferences:  prefs,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// UpdatePreferences updates the user preferences and refreshes the UpdatedAt timestamp.
func (u *User) UpdatePreferences(prefs UserPreferences) error {
	if u.IsDeleted() {
		return ErrUserAlreadyDeleted
	}
	if err := prefs.Validate(); err != nil {
		return err
	}
	u.Preferences = UserPreferences{BaseCurrency: strings.ToUpper(strings.TrimSpace(prefs.BaseCurrency))}
	u.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdatePasswordHash updates the password hash if valid and not deleted.
func (u *User) UpdatePasswordHash(newHash string) error {
	if u.IsDeleted() {
		return ErrUserAlreadyDeleted
	}
	clean := strings.TrimSpace(newHash)
	if clean == "" {
		return ErrEmptyPasswordHash
	}
	u.PasswordHash = clean
	u.UpdatedAt = time.Now().UTC()
	return nil
}

// Delete marks the user as soft-deleted.
func (u *User) Delete() error {
	if u.IsDeleted() {
		return ErrUserAlreadyDeleted
	}
	now := time.Now().UTC()
	u.DeletedAt = &now
	u.UpdatedAt = now
	return nil
}

// IsDeleted checks if the user has been soft-deleted.
func (u *User) IsDeleted() bool {
	return u.DeletedAt != nil
}

