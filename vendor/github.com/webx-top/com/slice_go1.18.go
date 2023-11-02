//go:build go1.18

package com

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
