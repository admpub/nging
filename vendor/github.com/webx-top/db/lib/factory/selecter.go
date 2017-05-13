package factory

import (
	"math"
	"math/rand"
	"sync/atomic"
)

type Selecter interface {
	Select(length int) int
}

type Random struct{}

func (r *Random) Select(length int) int {
	return rand.Intn(length - 1)
}

type RoundRobin struct {
	count uint64
}

func (r *RoundRobin) Select(length int) int {
	var c uint64
	if atomic.LoadUint64(&r.count) >= math.MaxUint64-1 {
		atomic.StoreUint64(&r.count, 0)
		c = 1
	} else {
		c = atomic.AddUint64(&r.count, 1)
	}
	return int(c % uint64(length))
}

func NewWeightedRoundRobin(length int) *weightedRoundRobin {
	return &weightedRoundRobin{
		weights:    make([]int32, length),
		lastSelect: -1,
		gcd:        2,
	}
}

type weightedRoundRobin struct {
	weights       []int32
	maxWeight     int32
	lastSelect    int32
	currentWeight int32
	gcd           int32
}

func (w *weightedRoundRobin) Select(length int) int {
	if w.gcd < 0 {
		w.gcd = MathNGCD(w.weights, length)
	}
	ln := int32(length)
	for {
		i := atomic.AddInt32(&w.lastSelect, 1) % ln
		atomic.StoreInt32(&w.lastSelect, i)
		if i == 0 {
			cw := atomic.AddInt32(&w.currentWeight, -w.gcd)
			if cw <= 0 {
				cw = w.maxWeight
				atomic.StoreInt32(&w.currentWeight, cw)
				if cw == 0 {
					return 0
				}
			} else {
				atomic.StoreInt32(&w.currentWeight, cw)
			}
		}
		j := int(i)
		if w.weights[j] >= atomic.LoadInt32(&w.currentWeight) {
			return j
		}
	}
}

func (w *weightedRoundRobin) SetWeight(index int, weight int32) *weightedRoundRobin {
	if index < len(w.weights) {
		w.weights[index] = weight
	} else {
		for j := len(w.weights); j <= index; j++ {
			w.weights = append(w.weights, 0)
		}
		w.weights[index] = weight
	}
	if w.maxWeight < weight {
		w.maxWeight = weight
	}
	w.gcd = -1
	return w
}

type SelectFirst struct {
}

func (s *SelectFirst) Select(length int) int {
	return 0
}
