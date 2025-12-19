package domain

import (
	"time"

	"github.com/google/uuid"
)

type Brand struct {
	Id        int32      `json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	Logo      string     `json:"logo"`
	CreatedBy uuid.UUID  `json:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedBy *uuid.UUID `json:"-"`
	UpdatedAt *time.Time `json:"-"`
}
