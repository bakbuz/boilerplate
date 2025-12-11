package dto

import (
	"codegen/internal/entity"

	"github.com/google/uuid"
)

type ProductSearchFilter struct {
	Id      *uuid.UUID
	BrandId int
	Name    string
	Skip    int16
	Take    int16
}

type ProductSearchResult struct {
	Total int
	Items []*entity.Product
}
