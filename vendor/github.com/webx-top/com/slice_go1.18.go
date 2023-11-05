//go:build go1.18

package com

import "sort"

type Number interface {
	~uint8 | ~int8 | ~uint16 | ~int16 | ~uint32 | ~int32 | ~uint | ~int | ~uint64 | ~int64 | ~float32 | ~float64
}

type Scalar interface {
	Number | ~bool | ~string
}

func SliceExtractCallback[T Scalar](parts []string, cb func(string) T, recv ...*T) {
	recvEndIndex := len(recv) - 1
	if recvEndIndex < 0 {
		return
	}
	for index, value := range parts {
		if index > recvEndIndex {
			break
		}
		*recv[index] = cb(value)
	}
}

type reverseSortIndex[T any] []T

func (s reverseSortIndex[T]) Len() int { return len(s) }
func (s reverseSortIndex[T]) Less(i, j int) bool {
	return j < i
}
func (s reverseSortIndex[T]) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func ReverseSortIndex[T any](values []T) []T {
	sort.Sort(reverseSortIndex[T](values))
	return values
}
