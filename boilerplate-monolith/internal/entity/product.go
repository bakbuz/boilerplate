package entity

import (
	"codegen/pkg/errx"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Product struct {
	Id            uuid.UUID  `json:"id"`
	BrandId       int        `json:"brandId"`
	Name          string     `json:"name"`
	Sku           *string    `json:"sku"`
	Summary       *string    `json:"summary"`
	Storyline     *string    `json:"storyline"`
	StockQuantity int        `json:"stockQuantity"`
	Price         float64    `json:"price"`
	Deleted       bool       `json:"-"`
	CreatedBy     uuid.UUID  `json:"-"`
	CreatedAt     time.Time  `json:"-"`
	UpdatedBy     *uuid.UUID `json:"-"`
	UpdatedAt     *time.Time `json:"-"`
	DeletedBy     *uuid.UUID `json:"-"`
	DeletedAt     *time.Time `json:"-"`
}

// Validate product entity for create/update operations
func (e *Product) Validate() error {
	if e == nil {
		return errx.ErrInvalidInput
	}

	// Sanitize strings
	e.Name = strings.TrimSpace(e.Name)
	if e.Sku != nil {
		trimmed := strings.TrimSpace(*e.Sku)
		e.Sku = &trimmed
	}
	if e.Summary != nil {
		trimmed := strings.TrimSpace(*e.Summary)
		e.Summary = &trimmed
	}
	if e.Storyline != nil {
		trimmed := strings.TrimSpace(*e.Storyline)
		e.Storyline = &trimmed
	}

	// Validate product name
	if e.Name == "" {
		return errors.New("product name is required")
	}

	if len(e.Name) > 255 {
		return errors.New("product name must not exceed 255 characters")
	}

	// Validate SKU if provided
	if e.Sku != nil && len(*e.Sku) > 100 {
		return errors.New("product SKU must not exceed 100 characters")
	}

	// Validate summary if provided
	if e.Summary != nil && len(*e.Summary) > 500 {
		return errors.New("product summary must not exceed 500 characters")
	}

	// Validate storyline if provided
	if e.Storyline != nil && len(*e.Storyline) > 2000 {
		return errors.New("product storyline must not exceed 2000 characters")
	}

	// Validate brand Id
	if e.BrandId <= 0 {
		return errors.New("valid brand Id is required")
	}

	// Validate stock quantity
	if e.StockQuantity < 0 {
		return errors.New("stock quantity cannot be negative")
	}

	// Validate price
	if e.Price < 0 {
		return errors.New("price cannot be negative")
	}

	return nil
}
