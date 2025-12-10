package entity

import "time"

type Product struct {
	ID          int64
	SKU         string
	Name        string
	Description string
	Price       float64
	Stock       int32
	CreatedAt   time.Time
}
