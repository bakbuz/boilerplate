package pkg

import (
	"codegen/internal/entity"

	"github.com/google/uuid"
)

type ProductFilterReq struct {
	Id   *uuid.UUID `query:"id"`
	Name string     `query:"name"`
	// SortCol *string `query:"sortCol"`
	// SortDir *byte   `query:"sortDir"` // 0 ASC, 1 DESC
	Page  int16 `query:"page"`
	Limit int16 `query:"limit"`
}

type ProductsResponse struct {
	Count    int               `json:"count"`
	Products []*entity.Product `json:"products"`
}

type ProductResponse struct {
	Product *entity.Product `json:"product"`
}

type ProductCreateUpdateReq struct {
	BrandId       int     `json:"brandId" validate:"required"`
	Name          string  `json:"name" validate:"required,max=100"`
	Sku           string  `json:"sku" validate:"max=50"`
	Summary       string  `json:"summary" validate:"max=500"`
	Storyline     string  `json:"storyline"`
	StockQuantity int     `json:"stockQuantity" validate:"required"`
	Price         float64 `json:"price" validate:"required"`
}
