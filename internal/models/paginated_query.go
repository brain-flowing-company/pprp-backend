package models

import (
	"fmt"

	"gorm.io/gorm"
)

type paginator struct {
	offset int
	limit  int
}

func NewPaginatedQuery(page int, limit int) *paginator {
	return &paginator{
		offset: (page - 1) * limit,
		limit:  limit,
	}
}

func (p *paginator) PaginatedQuery(db *gorm.DB) *gorm.DB {
	return db.Offset(p.offset).Limit(p.limit)
}

func (p *paginator) PaginatedSQL() string {
	return fmt.Sprintf("LIMIT %d OFFSET %d", p.limit, p.offset)
}
