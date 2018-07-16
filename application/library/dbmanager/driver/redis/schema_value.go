package redis

import (
	"strconv"
	"strings"
)

func NewValueInfo(r *Redis,k string) &ValueInfo{
	return &ValueInfo{
		r:r,
		k:k,
	}
}

type ValueInfo struct {
	sizeNum  int
	sizeUnit string
	encoding string
	ttl      int
	typeName string
	k      string
	r    *Redis
}

func (v *ValueInfo) Size() {

}
