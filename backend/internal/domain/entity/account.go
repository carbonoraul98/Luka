package entity

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrEmptyAccountID         = errors.New("account ID cannot be empty")
	ErrEmptyAccountUserID     = errors.New("user ID cannot be empty")
	ErrEmptyAccountName       = errors.New("account name cannot be empty")
	ErrInvalidAccountType     = errors.New("invalid account type")
	ErrInvalidAccountCategory = errors.New("invalid account category")
	ErrParentAccountRequired  = errors.New("parent account ID is required for savings pockets")
	ErrParentAccountForbidden = errors.New("parent account ID can only be set for savings pockets")
	ErrNegativeBalance        = errors.New("account balance cannot be negative")
	ErrNegativeTargetAmount   = errors.New("target amount must be greater than zero")
	ErrInsufficientFunds      = errors.New("insufficient funds for withdrawal")
	ErrInvalidAmount          = errors.New("amount must be greater than zero")
	ErrAccountAlreadyDeleted  = errors.New("account is already deleted")
)

// AccountCategory represents the classification of financial accounts.
type AccountCategory string

const (
	CategoryCash          AccountCategory = "CASH"
	CategoryBank          AccountCategory = "BANK"
	CategoryWallet        AccountCategory = "WALLET"
	CategorySavingsPocket AccountCategory = "SAVINGS_POCKET"
)

// IsValid checks if an account category is supported.
func (c AccountCategory) IsValid() bool {
	switch c {
	case CategoryCash, CategoryBank, CategoryWallet, CategorySavingsPocket:
		return true
	default:
		return false
	}
}

// IsValid checks if an account type is supported.
func (t AccountType) IsValid() bool {
	switch t {
	case AccountPersonal, AccountWorkingCapital:
		return true
	default:
		return false
	}
}

// Account represents a monetary holding entity or a savings pocket.
type Account struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id"`
	ParentAccountID *string         `json:"parent_account_id,omitempty"`
	Name            string          `json:"name"`
	Type            AccountType     `json:"type"`
	Category        AccountCategory `json:"category"`
	Currency        string          `json:"currency"`
	Balance         float64         `json:"balance"`
	TargetAmount    *float64        `json:"target_amount,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
	DeletedAt       *time.Time      `json:"deleted_at,omitempty"`
}

// NewAccount creates and validates a standard bank, cash, or wallet account.
func NewAccount(id, userID, name string, accType AccountType, category AccountCategory, currency string, initialBalance float64) (*Account, error) {
	if category == CategorySavingsPocket {
		return nil, ErrParentAccountRequired
	}
	return createAccount(id, userID, nil, name, accType, category, currency, initialBalance, nil)
}

// NewSavingsPocket creates a savings sub-account linked to a parent account.
func NewSavingsPocket(id, userID, parentAccountID, name string, parentType AccountType, currency string, targetAmount *float64) (*Account, error) {
	cleanParentID := strings.TrimSpace(parentAccountID)
	if cleanParentID == "" {
		return nil, ErrParentAccountRequired
	}

	if targetAmount != nil && *targetAmount <= 0 {
		return nil, ErrNegativeTargetAmount
	}

	return createAccount(id, userID, &cleanParentID, name, parentType, CategorySavingsPocket, currency, 0.0, targetAmount)
}

func createAccount(
	id, userID string,
	parentID *string,
	name string,
	accType AccountType,
	category AccountCategory,
	currency string,
	balance float64,
	targetAmount *float64,
) (*Account, error) {
	cleanID := strings.TrimSpace(id)
	if cleanID == "" {
		return nil, ErrEmptyAccountID
	}

	cleanUserID := strings.TrimSpace(userID)
	if cleanUserID == "" {
		return nil, ErrEmptyAccountUserID
	}

	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return nil, ErrEmptyAccountName
	}

	if !accType.IsValid() {
		return nil, ErrInvalidAccountType
	}

	if !category.IsValid() {
		return nil, ErrInvalidAccountCategory
	}

	if category != CategorySavingsPocket && parentID != nil {
		return nil, ErrParentAccountForbidden
	}

	cleanCurrency := strings.ToUpper(strings.TrimSpace(currency))
	if _, ok := SupportedCurrencies[cleanCurrency]; !ok {
		return nil, ErrUnsupportedCurrency
	}

	if balance < 0 {
		return nil, ErrNegativeBalance
	}

	now := time.Now().UTC()
	return &Account{
		ID:              cleanID,
		UserID:          cleanUserID,
		ParentAccountID: parentID,
		Name:            cleanName,
		Type:            accType,
		Category:        category,
		Currency:        cleanCurrency,
		Balance:         balance,
		TargetAmount:    targetAmount,
		CreatedAt:       now,
		UpdatedAt:       now,
	}, nil
}

// Deposit increments the account balance by amount.
func (a *Account) Deposit(amount float64) error {
	if a.IsDeleted() {
		return ErrAccountAlreadyDeleted
	}
	if amount <= 0 {
		return ErrInvalidAmount
	}
	a.Balance += amount
	a.UpdatedAt = time.Now().UTC()
	return nil
}

// Withdraw decrements the account balance if funds are sufficient.
func (a *Account) Withdraw(amount float64) error {
	if a.IsDeleted() {
		return ErrAccountAlreadyDeleted
	}
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if a.Balance < amount {
		return ErrInsufficientFunds
	}
	a.Balance -= amount
	a.UpdatedAt = time.Now().UTC()
	return nil
}

// IsSavingsPocket indicates whether this account is a pocket.
func (a *Account) IsSavingsPocket() bool {
	return a.Category == CategorySavingsPocket && a.ParentAccountID != nil
}

// Delete soft-deletes the account.
func (a *Account) Delete() error {
	if a.IsDeleted() {
		return ErrAccountAlreadyDeleted
	}
	now := time.Now().UTC()
	a.DeletedAt = &now
	a.UpdatedAt = now
	return nil
}

// IsDeleted checks if the account has been soft-deleted.
func (a *Account) IsDeleted() bool {
	return a.DeletedAt != nil
}

