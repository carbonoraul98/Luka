package entity

import (
	"errors"
	"time"
)

var (
	ErrZeroQuantity      = errors.New("quantity must be greater than zero")
	ErrNegativeCost      = errors.New("costs cannot be negative")
	ErrNoItemsInBatch    = errors.New("batch must have at least one product item")
	ErrAllocationMismatch = errors.New("cost allocation error: total weight or price sum is zero")
)

// ItemStatus defines the lifecycle of an inventory item.
type ItemStatus string

const (
	StatusInTransit ItemStatus = "IN_TRANSIT"
	StatusInStock   ItemStatus = "IN_STOCK"
	StatusSold      ItemStatus = "SOLD"
)

// AccountType distinguishes working capital from personal funds.
type AccountType string

const (
	AccountPersonal       AccountType = "PERSONAL"
	AccountWorkingCapital AccountType = "WORKING_CAPITAL"
)

// BatchItem represents a specific product line within an import batch.
type BatchItem struct {
	ID              string  `json:"id"`
	ProductID       string  `json:"product_id"`
	Name            string  `json:"name"`
	Quantity        int     `json:"quantity"`
	SupplierPrice   float64 `json:"supplier_price"`   // Unit price in foreign currency or USD
	WeightKg        float64 `json:"weight_kg"`        // Useful for weight-based shipping allocation
	CalculatedLandedUnitCost float64 `json:"calculated_landed_unit_cost"` // Landed Cost per unit in local currency
}

// ImportBatch aggregates products bought abroad and all additional landed costs.
type ImportBatch struct {
	ID               string      `json:"id"`
	Title            string      `json:"title"`
	OriginPlatform   string      `json:"origin_platform"` // e.g. "AliExpress"
	Currency         string      `json:"currency"`        // e.g. "USD"
	ExchangeRate     float64     `json:"exchange_rate"`   // Local currency per 1 foreign unit
	ShippingFee      float64     `json:"shipping_fee"`    // Total international freight
	CustomsDutyFee   float64     `json:"customs_duty_fee"`// Tariffs / Taxes
	OtherFees        float64     `json:"other_fees"`      // Local courier, handling, bank fees
	Items            []BatchItem `json:"items"`
	Status           ItemStatus  `json:"status"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	DeletedAt        *time.Time  `json:"deleted_at,omitempty"`
}

// CalculateLandedCosts prorates shipping, customs and fees across all batch items.
// Prorating method: weighted by total item value (SupplierPrice * Quantity).
func (b *ImportBatch) CalculateLandedCosts() error {
	if len(b.Items) == 0 {
		return ErrNoItemsInBatch
	}

	totalProductValue := 0.0
	for _, item := range b.Items {
		if item.Quantity <= 0 {
			return ErrZeroQuantity
		}
		if item.SupplierPrice < 0 {
			return ErrNegativeCost
		}
		totalProductValue += item.SupplierPrice * float64(item.Quantity)
	}

	if totalProductValue == 0 {
		return ErrAllocationMismatch
	}

	totalAdditionalFees := b.ShippingFee + b.CustomsDutyFee + b.OtherFees

	for i := range b.Items {
		itemValue := b.Items[i].SupplierPrice * float64(b.Items[i].Quantity)
		proportion := itemValue / totalProductValue

		// Allocated extra cost for this line item
		allocatedFees := totalAdditionalFees * proportion
		totalLineCost := itemValue + allocatedFees

		// Unit landed cost in foreign currency converted to local currency
		unitLandedCostForeign := totalLineCost / float64(b.Items[i].Quantity)
		rate := b.ExchangeRate
		if rate <= 0 {
			rate = 1.0
		}
		b.Items[i].CalculatedLandedUnitCost = unitLandedCostForeign * rate
	}

	return nil
}
