package domain

type Identifier int64

func Identifiers(values []int64) []Identifier {
	ids := make([]Identifier, 0, len(values))
	for _, value := range values {
		ids = append(ids, Identifier(value))
	}
	return ids
}
