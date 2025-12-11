package entity

import "time"

type Brand struct {
	Id        int32      `json:"id"`
	Name      string     `json:"name"`
	Slug      string     `json:"slug"`
	Logo      string     `json:"logo"`
	CreatedBy int        `json:"-"`
	CreatedAt time.Time  `json:"-"`
	UpdatedBy *int       `json:"-"`
	UpdatedAt *time.Time `json:"-"`
}
