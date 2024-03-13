package models

import (
	"gorm.io/gorm"
)

type PaginatedQuery struct {
	Offset int
	Limit  int
}

func NewPaginatedQuery(page int, limit int) *PaginatedQuery {
	return &PaginatedQuery{
		Offset: (page - 1) * limit,
		Limit:  limit,
	}
}

func (p *PaginatedQuery) PaginatedQuery(db *gorm.DB) *gorm.DB {
	return db.Offset(p.Offset).Limit(p.Limit)
}
