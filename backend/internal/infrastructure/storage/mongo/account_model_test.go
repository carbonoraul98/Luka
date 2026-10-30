package mongo_test

import (
	"testing"
	"time"

	"github.com/raulcarbono/luka/backend/internal/domain/entity"
	mongostorage "github.com/raulcarbono/luka/backend/internal/infrastructure/storage/mongo"
)

func TestAccountModel_Mapping(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Millisecond)
	target := 500000.0
	parentID := "acc-parent-123"

	pocket := &entity.Account{
		ID:              "pocket-uuid-1",
		UserID:          "usr-uuid-1",
		ParentAccountID: &parentID,
		Name:            "Ahorro Navidad",
		Type:            entity.AccountPersonal,
		Category:        entity.CategorySavingsPocket,
		Currency:        "COP",
		Balance:         150000.0,
		TargetAmount:    &target,
		CreatedAt:       now,
		UpdatedAt:       now,
		DeletedAt:       nil,
	}

	doc := mongostorage.FromAccountDomain(pocket)

	if doc.ID != pocket.ID {
		t.Errorf("got doc ID %q, want %q", doc.ID, pocket.ID)
	}
	if doc.UserID != pocket.UserID {
		t.Errorf("got doc UserID %q, want %q", doc.UserID, pocket.UserID)
	}
	if doc.ParentAccountID == nil || *doc.ParentAccountID != parentID {
		t.Errorf("got doc ParentAccountID %v, want %q", doc.ParentAccountID, parentID)
	}
	if doc.Category != string(entity.CategorySavingsPocket) {
		t.Errorf("got doc Category %q, want %q", doc.Category, entity.CategorySavingsPocket)
	}
	if doc.TargetAmount == nil || *doc.TargetAmount != target {
		t.Errorf("got doc TargetAmount %v, want %f", doc.TargetAmount, target)
	}

	restored := doc.ToDomain()

	if restored.ID != pocket.ID {
		t.Errorf("got restored ID %q, want %q", restored.ID, pocket.ID)
	}
	if !restored.IsSavingsPocket() {
		t.Error("restored entity should be recognized as savings pocket")
	}
	if restored.Balance != pocket.Balance {
		t.Errorf("got restored Balance %f, want %f", restored.Balance, pocket.Balance)
	}
}

