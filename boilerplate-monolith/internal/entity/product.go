package entity

import (
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
