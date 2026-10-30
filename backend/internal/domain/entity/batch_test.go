package entity_test

import (
	"testing"

	"github.com/raulcarbono/luka/backend/internal/domain/entity"
)

func TestCalculateLandedCosts(t *testing.T) {
	batch := entity.ImportBatch{
		ID:             "batch-1",
		Title:          "AliExpress Tech Order #1",
		OriginPlatform: "AliExpress",
		Currency:       "USD",
		ExchangeRate:   4000.0, // e.g. COP to USD
		ShippingFee:    30.0,
		CustomsDutyFee: 20.0,
		OtherFees:      0.0,
		Items: []entity.BatchItem{
			{
				ID:            "item-1",
				Name:          "Smartwatch X",
				Quantity:      10,
				SupplierPrice: 15.0, // Total: 150 USD (75% of product spend)
			},
			{
				ID:            "item-2",
				Name:          "Earbuds Pro",
				Quantity:      5,
				SupplierPrice: 10.0, // Total: 50 USD (25% of product spend)
			},
		},
	}

	err := batch.CalculateLandedCosts()
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// Total product spend: 200 USD
	// Total fees: 50 USD (25% overhead on top of product spend)
	// Smartwatch landed foreign unit: 15 * 1.25 = 18.75 USD -> in local: 18.75 * 4000 = 75,000
	expectedSmartwatchCost := 75000.0
	actualSmartwatchCost := batch.Items[0].CalculatedLandedUnitCost
	if actualSmartwatchCost != expectedSmartwatchCost {
		t.Errorf("expected Smartwatch unit landed cost %f, got %f", expectedSmartwatchCost, actualSmartwatchCost)
	}

	// Earbuds landed foreign unit: 10 * 1.25 = 12.5 USD -> in local: 12.5 * 4000 = 50,000
	expectedEarbudsCost := 50000.0
	actualEarbudsCost := batch.Items[1].CalculatedLandedUnitCost
	if actualEarbudsCost != expectedEarbudsCost {
		t.Errorf("expected Earbuds unit landed cost %f, got %f", expectedEarbudsCost, actualEarbudsCost)
	}
}
