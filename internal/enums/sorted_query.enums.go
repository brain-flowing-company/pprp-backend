package enums

type SortDirection string

const (
	ASC  SortDirection = "asc"
	DESC SortDirection = "desc"
)

func ParseSortDirection(dir string) (SortDirection, bool) {
	val, ok := map[string]SortDirection{
		"asc":  ASC,
		"desc": DESC,
	}[dir]
	return val, ok
}
