package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/brain-flowing-company/pprp-backend/internal/enums"
)

type FilteredQuery struct {
	items  []string
	mapper map[string]string
}

func NewFilteredQuery(model interface{}) *FilteredQuery {
	s := &FilteredQuery{mapper: map[string]string{}}
	t := reflect.TypeOf(model)

	parents := NewStack[string]()
	s.assignMapper(parents, t)

	return s
}

func (s *FilteredQuery) assignMapper(parents *Stack[string], t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		json := f.Tag.Get("json")
		sortmap := f.Tag.Get("filtermapper")

		if json == "-" || sortmap == "-" {
			continue
		}

		if f.Type.Kind() == reflect.Struct {
			if len(json) > 0 {
				parents.Push(json)
			}

			s.assignMapper(parents, f.Type)

			if len(json) > 0 {
				parents.Pop()
			}
		}

		if len(json) == 0 || len(sortmap) == 0 {
			continue
		}

		parents.Push(json)
		key := strings.Join(parents.Seek(), ".")
		s.mapper[key] = sortmap
		parents.Pop()
	}
}

func (s *FilteredQuery) ParseQuery(query string) error {
	if len(query) == 0 {
		return nil
	}

	filters := strings.Split(query, ",")

	for _, filter := range filters {
		ob, cb, cl := -1, -1, -1
		for i, ch := range filter {
			switch ch {
			case '[':
				ob = i
			case ']':
				cb = i
			case ':':
				cl = i
			}
		}

		if ob == -1 || cb == -1 || cl == -1 {
			return fmt.Errorf("%s invalid format, <field>[<opt>]:<value>", filter)
		}

		fld := filter[:ob]
		field, ok := s.mapper[fld]
		if !ok {
			return fmt.Errorf("'%s' is not a valid sort key", fld)
		}

		opt := filter[ob+1 : cb]
		operation, ok := enums.ParseFilterOperation(opt)
		if !ok {
			return errors.New("sort direction can only be 'lte' or 'gte'")
		}

		var value float32
		_, err := fmt.Sscanf(filter[cl+1:], "%f", &value)
		if err != nil {
			return err
		}

		s.items = append(s.items, fmt.Sprintf("%s %s %f", field, operation, value))
	}

	return nil
}

func (s *FilteredQuery) Map(key string, value string) {
	s.mapper[key] = value
}

func (s *FilteredQuery) FilteredSQL() string {
	if len(s.items) > 0 {
		return strings.Join(s.items, " AND ")
	}
	return "TRUE"
}
