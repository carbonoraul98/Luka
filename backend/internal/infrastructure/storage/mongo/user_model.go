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

const UsersCollection = "users"

// UserPreferencesDocument represents the nested preferences document in MongoDB.
type UserPreferencesDocument struct {
	BaseCurrency string `bson:"base_currency"`
}

// UserDocument defines the physical MongoDB document schema.
type UserDocument struct {
	ID           string                  `bson:"_id"`
	Email        string                  `bson:"email"`
	PasswordHash string                  `bson:"password_hash"`
	Preferences  UserPreferencesDocument `bson:"preferences"`
	CreatedAt    time.Time               `bson:"created_at"`
	UpdatedAt    time.Time               `bson:"updated_at"`
	DeletedAt    *time.Time              `bson:"deleted_at,omitempty"`
}

// FromDomain maps a domain User entity to a MongoDB UserDocument.
func FromDomain(u *entity.User) UserDocument {
	return UserDocument{
		ID:           u.ID,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Preferences: UserPreferencesDocument{
			BaseCurrency: u.Preferences.BaseCurrency,
		},
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}
}

// ToDomain maps a MongoDB UserDocument back to the domain User entity.
func (d UserDocument) ToDomain() *entity.User {
	return &entity.User{
		ID:           d.ID,
		Email:        d.Email,
		PasswordHash: d.PasswordHash,
		Preferences: entity.UserPreferences{
			BaseCurrency: d.Preferences.BaseCurrency,
		},
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		DeletedAt: d.DeletedAt,
	}
}

// EnsureUserIndexes creates the required unique and lookup indexes for the users collection.
func EnsureUserIndexes(ctx context.Context, coll *mongo.Collection) error {
	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_users_email_unique"),
		},
		{
			Keys:    bson.D{{Key: "deleted_at", Value: 1}},
			Options: options.Index().SetName("idx_users_deleted_at"),
		},
	}

	_, err := coll.Indexes().CreateMany(ctx, models)
	if err != nil {
		return fmt.Errorf("create user indexes: %w", err)
	}
	return nil
}

