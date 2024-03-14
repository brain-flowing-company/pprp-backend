package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"gorm.io/gorm"
)

type sorter struct {
	Field     string
	Direction enums.SortDirection
	mapper    map[string]string
}

func NewSortedQuery(t reflect.Type) *sorter {
	s := &sorter{mapper: map[string]string{}}
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		json := f.Tag.Get("json")
		sortmap := f.Tag.Get("sortmapper")

		if len(json) > 0 && len(sortmap) > 0 {
			s.mapper[json] = sortmap
		}
	}

	return s
}

func (s *sorter) Parse(query string) error {
	pairs := strings.Split(query, ":")

	field, ok := s.mapper[pairs[0]]
	if !ok {
		return errors.New("invalid sort field")
	}

	direction, ok := enums.ParseSortDirection(pairs[1])
	if !ok {
		return errors.New("sort direction can only be 'asc' or 'desc'")
	}

	s.Field = field
	s.Direction = direction

	return nil
}

func (s *sorter) Map(key string, value string) {
	s.mapper[key] = value
}

func (s *sorter) SortedQuery(db *gorm.DB) *gorm.DB {
	return db.Order(fmt.Sprintf("%s %s", s.Field, s.Direction))
}
