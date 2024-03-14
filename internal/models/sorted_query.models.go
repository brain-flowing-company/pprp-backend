package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/brain-flowing-company/pprp-backend/internal/enums"
	"gorm.io/gorm"
)

type SortedQuery struct {
	Field     string
	Direction enums.SortDirection
	mapper    map[string]string
}

func NewSortedQuery(model interface{}) *SortedQuery {
	s := &SortedQuery{mapper: map[string]string{}}
	t := reflect.TypeOf(model)
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

func (s *SortedQuery) ParseQuery(query string) error {
	if len(query) == 0 {
		return nil
	}

	pairs := strings.Split(query, ":")

	if len(pairs) < 2 {
		return errors.New("too few sorting arguments")
	}

	field, ok := s.mapper[pairs[0]]
	if !ok {
		return fmt.Errorf("'%s' is not a valid sort key", pairs[0])
	}

	direction, ok := enums.ParseSortDirection(pairs[1])
	if !ok {
		return errors.New("sort direction can only be 'asc' or 'desc'")
	}

	s.Field = field
	s.Direction = direction

	return nil
}

func (s *SortedQuery) Map(key string, value string) {
	s.mapper[key] = value
}

func (s *SortedQuery) SortedQuery(db *gorm.DB) *gorm.DB {
	return db.Order(fmt.Sprintf("%s %s", s.Field, s.Direction))
}
