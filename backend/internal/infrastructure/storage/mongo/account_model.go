package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/raulcarbono/luka/backend/internal/domain/entity"
)

const AccountsCollection = "accounts"

// AccountDocument defines the physical MongoDB document schema for accounts.
type AccountDocument struct {
	ID              string     `bson:"_id"`
	UserID          string     `bson:"user_id"`
	ParentAccountID *string    `bson:"parent_account_id,omitempty"`
	Name            string     `bson:"name"`
	Type            string     `bson:"type"`
	Category        string     `bson:"category"`
	Currency        string     `bson:"currency"`
	Balance         float64    `bson:"balance"`
	TargetAmount    *float64   `bson:"target_amount,omitempty"`
	CreatedAt       time.Time  `bson:"created_at"`
	UpdatedAt       time.Time  `bson:"updated_at"`
	DeletedAt       *time.Time `bson:"deleted_at,omitempty"`
}

// FromAccountDomain maps a domain Account entity to a MongoDB AccountDocument.
func FromAccountDomain(a *entity.Account) AccountDocument {
	return AccountDocument{
		ID:              a.ID,
		UserID:          a.UserID,
		ParentAccountID: a.ParentAccountID,
		Name:            a.Name,
		Type:            string(a.Type),
		Category:        string(a.Category),
		Currency:        a.Currency,
		Balance:         a.Balance,
		TargetAmount:    a.TargetAmount,
		CreatedAt:       a.CreatedAt,
		UpdatedAt:       a.UpdatedAt,
		DeletedAt:       a.DeletedAt,
	}
}

// ToDomain maps a MongoDB AccountDocument back to the domain Account entity.
func (d AccountDocument) ToDomain() *entity.Account {
	return &entity.Account{
		ID:              d.ID,
		UserID:          d.UserID,
		ParentAccountID: d.ParentAccountID,
		Name:            d.Name,
		Type:            entity.AccountType(d.Type),
		Category:        entity.AccountCategory(d.Category),
		Currency:        d.Currency,
		Balance:         d.Balance,
		TargetAmount:    d.TargetAmount,
		CreatedAt:       d.CreatedAt,
		UpdatedAt:       d.UpdatedAt,
		DeletedAt:       d.DeletedAt,
	}
}

// EnsureAccountIndexes creates the performance and constraint indexes for the accounts collection.
func EnsureAccountIndexes(ctx context.Context, coll *mongo.Collection) error {
	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "name", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_accounts_user_name_unique"),
		},
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "type", Value: 1}},
			Options: options.Index().SetName("idx_accounts_user_type"),
		},
		{
			Keys:    bson.D{{Key: "parent_account_id", Value: 1}},
			Options: options.Index().SetName("idx_accounts_parent_id"),
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("idx_accounts_deleted_at"),
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("create account indexes: %w", err)
	}
	return nil
}

