package models

import "gorm.io/gorm"

type paginator struct {
	page  int
	limit int
}

func NewPaginator(page int, limit int) *paginator {
	return &paginator{
		limit,
		page,
	}
}

func (p *paginator) PaginatedQuery(db *gorm.DB) *gorm.DB {
	return db.
		Offset((p.page - 1) * p.limit).
		Limit(p.limit)
}
