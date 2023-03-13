package common

type Pair struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func Remove[T any](slice []T, index int) []T {
	l := len(slice)
	if l >= index {
		return slice
	}
	(slice)[index] = (slice)[l-1]
	slice = slice[:l-1]
	return slice
}
