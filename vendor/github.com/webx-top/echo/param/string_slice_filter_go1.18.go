//go:build go1.18

package param

type Number interface {
	~uint8 | ~int8 | ~uint16 | ~int16 | ~uint32 | ~int32 | ~uint | ~int | ~uint64 | ~int64 | ~float32 | ~float64
}

type Scalar interface {
	Number | ~bool | ~string
}

func IsGreaterThanZeroElement[T Number](_ int, v T) bool {
	return v > 0
}

func Unique[T comparable](p []T) []T {
	record := map[T]struct{}{}
	var result []T
	for _, s := range p {
		if _, ok := record[s]; !ok {
			record[s] = struct{}{}
			result = append(result, s)
		}
	}
	return result
}
