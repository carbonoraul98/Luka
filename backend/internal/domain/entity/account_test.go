package entity_test

import (
	"errors"
	"testing"

	"github.com/raulcarbono/luka/backend/internal/domain/entity"
)

func TestNewAccount(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		userID         string
		accName        string
		accType        entity.AccountType
		category       entity.AccountCategory
		currency       string
		initialBalance float64
		wantErr        error
	}{
		{
			name:           "valid cash account",
			id:             "acc-1",
			userID:         "usr-1",
			accName:        "Petty Cash",
			accType:        entity.AccountWorkingCapital,
			category:       entity.CategoryCash,
			currency:       "COP",
			initialBalance: 50000.0,
			wantErr:        nil,
		},
		{
			name:           "valid bank account",
			id:             "acc-2",
			userID:         "usr-1",
			accName:        "Main Bank Account",
			accType:        entity.AccountPersonal,
			category:       entity.CategoryBank,
			currency:       "USD",
			initialBalance: 0.0,
			wantErr:        nil,
		},
		{
			name:           "empty id",
			id:             "  ",
			userID:         "usr-1",
			accName:        "Main",
			accType:        entity.AccountPersonal,
			category:       entity.CategoryBank,
			currency:       "USD",
			initialBalance: 0,
			wantErr:        entity.ErrEmptyAccountID,
		},
		{
			name:           "empty user id",
			id:             "acc-3",
			userID:         " ",
			accName:        "Main",
			accType:        entity.AccountPersonal,
			category:       entity.CategoryBank,
			currency:       "USD",
			initialBalance: 0,
			wantErr:        entity.ErrEmptyAccountUserID,
		},
		{
			name:           "empty name",
			id:             "acc-4",
			userID:         "usr-1",
			accName:        "",
			accType:        entity.AccountPersonal,
			category:       entity.CategoryBank,
			currency:       "USD",
			initialBalance: 0,
			wantErr:        entity.ErrEmptyAccountName,
		},
		{
			name:           "negative initial balance",
			id:             "acc-5",
			userID:         "usr-1",
			accName:        "Bank",
			accType:        entity.AccountPersonal,
			category:       entity.CategoryBank,
			currency:       "USD",
			initialBalance: -100,
			wantErr:        entity.ErrNegativeBalance,
		},
		{
			name:           "cannot create pocket directly through NewAccount",
			id:             "acc-6",
			userID:         "usr-1",
			accName:        "Pocket",
			accType:        entity.AccountPersonal,
			category:       entity.CategorySavingsPocket,
			currency:       "USD",
			initialBalance: 0,
			wantErr:        entity.ErrParentAccountRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, err := entity.NewAccount(tt.id, tt.userID, tt.accName, tt.accType, tt.category, tt.currency, tt.initialBalance)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("expected error %v, got %v", tt.wantErr, err)
				}
				if acc != nil {
					t.Fatalf("expected nil account, got %+v", acc)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if acc.Balance != tt.initialBalance {
				t.Errorf("got balance %f, want %f", acc.Balance, tt.initialBalance)
			}
			if acc.IsSavingsPocket() {
				t.Error("standard account should not be a savings pocket")
			}
		})
	}
}

func TestNewSavingsPocket(t *testing.T) {
	target := 1000.0
	invalidTarget := -50.0

	t.Run("valid savings pocket", func(t *testing.T) {
		pocket, err := entity.NewSavingsPocket("pkt-1", "usr-1", "acc-parent-1", "Emergency Fund", entity.AccountPersonal, "USD", &target)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !pocket.IsSavingsPocket() {
			t.Error("account should be a savings pocket")
		}
		if pocket.ParentAccountID == nil || *pocket.ParentAccountID != "acc-parent-1" {
			t.Errorf("got parent ID %v, want acc-parent-1", pocket.ParentAccountID)
		}
		if pocket.Balance != 0.0 {
			t.Errorf("pocket initial balance must be 0, got %f", pocket.Balance)
		}
	})

	t.Run("empty parent account ID", func(t *testing.T) {
		_, err := entity.NewSavingsPocket("pkt-2", "usr-1", "  ", "Emergency Fund", entity.AccountPersonal, "USD", nil)
		if !errors.Is(err, entity.ErrParentAccountRequired) {
			t.Fatalf("expected ErrParentAccountRequired, got %v", err)
		}
	})

	t.Run("negative target amount", func(t *testing.T) {
		_, err := entity.NewSavingsPocket("pkt-3", "usr-1", "acc-parent-1", "Emergency Fund", entity.AccountPersonal, "USD", &invalidTarget)
		if !errors.Is(err, entity.ErrNegativeTargetAmount) {
			t.Fatalf("expected ErrNegativeTargetAmount, got %v", err)
		}
	})
}

func TestAccount_DepositAndWithdraw(t *testing.T) {
	acc, err := entity.NewAccount("acc-1", "usr-1", "Bank", entity.AccountPersonal, entity.CategoryBank, "USD", 100.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("valid deposit", func(t *testing.T) {
		if err := acc.Deposit(50.0); err != nil {
			t.Fatalf("deposit failed: %v", err)
		}
		if acc.Balance != 150.0 {
			t.Errorf("expected 150.0, got %f", acc.Balance)
		}
	})

	t.Run("deposit invalid amount", func(t *testing.T) {
		if err := acc.Deposit(-10.0); !errors.Is(err, entity.ErrInvalidAmount) {
			t.Fatalf("expected ErrInvalidAmount, got %v", err)
		}
	})

	t.Run("valid withdrawal", func(t *testing.T) {
		if err := acc.Withdraw(50.0); err != nil {
			t.Fatalf("withdraw failed: %v", err)
		}
		if acc.Balance != 100.0 {
			t.Errorf("expected 100.0, got %f", acc.Balance)
		}
	})

	t.Run("insufficient funds withdrawal", func(t *testing.T) {
		if err := acc.Withdraw(500.0); !errors.Is(err, entity.ErrInsufficientFunds) {
			t.Fatalf("expected ErrInsufficientFunds, got %v", err)
		}
	})
}

