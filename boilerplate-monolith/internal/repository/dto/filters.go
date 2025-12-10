package dto

import "github.com/google/uuid"

type ProductFilter struct {
	Id   *uuid.UUID
	Name string
	Skip int16
	Take int16
}
