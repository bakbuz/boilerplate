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
	CreatedBy     int        `json:"-"`
	CreatedAt     time.Time  `json:"-"`
	UpdatedBy     *int       `json:"-"`
	UpdatedAt     *time.Time `json:"-"`
	DeletedBy     *int       `json:"-"`
	DeletedAt     *time.Time `json:"-"`
}
