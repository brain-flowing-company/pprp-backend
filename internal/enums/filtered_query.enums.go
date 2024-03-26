package enums

type FilterOperation string

const (
	GTE FilterOperation = ">="
	LTE FilterOperation = "<="
	EQL FilterOperation = "="
)

func ParseFilterOperation(dir string) (FilterOperation, bool) {
	val, ok := map[string]FilterOperation{
		"gte": GTE,
		"lte": LTE,
		"eql": EQL,
	}[dir]
	return val, ok
}
