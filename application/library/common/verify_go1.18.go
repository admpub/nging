//go:build go1.18

package common

func If[T any](condition bool, yesValue T, noValue T) T {
	if condition {
		return yesValue
	}
	return noValue
}
