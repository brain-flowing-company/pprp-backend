package enums

type FilterOperation string

const (
	GTE FilterOperation = ">="
	LTE FilterOperation = "<="
)

func ParseFilterOperation(dir string) (FilterOperation, bool) {
	val, ok := map[string]FilterOperation{
		"gte": GTE,
		"lte": LTE,
	}[dir]
	return val, ok
}
