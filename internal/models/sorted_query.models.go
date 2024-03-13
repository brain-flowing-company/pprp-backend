package models

import "strings"

type sorter struct {
	field     string
	direction string
}

func NewSortedQuery(sort string) *sorter {
	pairs := strings.Split(sort, ":")
	return &sorter{
		field:     pairs[0],
		direction: pairs[1],
	}
}
