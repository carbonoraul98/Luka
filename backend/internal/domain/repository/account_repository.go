package repository

import (
	"context"
	"errors"

	"github.com/raulcarbono/luka/backend/internal/domain/entity"
)

var (
	ErrAccountNotFound          = errors.New("account not found")
	ErrAccountNameAlreadyExists = errors.New("account with this name already exists for this user")
)

// AccountRepository defines the persistence contract for Account entities.
type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) error
	FindByID(ctx context.Context, id string) (*entity.Account, error)
	FindByUserID(ctx context.Context, userID string) ([]*entity.Account, error)
	FindPocketsByParentID(ctx context.Context, parentID string) ([]*entity.Account, error)
	Update(ctx context.Context, account *entity.Account) error
}

