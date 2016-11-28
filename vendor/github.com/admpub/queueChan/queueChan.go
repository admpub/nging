// queueChan.go
package queueChan

import (
	"errors"
	"sync"
)

// QueueChan holds the queue within a ring buffer

type QueueChan struct {
	sync.RWMutex
	dynamic bool
	q       chan interface{}
	out     chan interface{}
	Empty   chan struct{}
}

func New(length ...int) QueueChan {
	capacity := 16
	if len(length) == 1 && length[0] > capacity {
		capacity = 1
		for ll := length[0]; ll > 0; ll >>= 1 {
			capacity <<= 1
		}
	}
	qu := &QueueChan{
		q:     make(chan interface{}, capacity),
		out:   make(chan interface{}, 1),
		Empty: make(chan struct{}),
	}
	return *qu
}

func (qu *QueueChan) New(length ...int) *QueueChan {
	capacity := 16
	if len(length) == 1 && length[0] > capacity {
		capacity = 1
		for ll := length[0]; ll > 0; ll >>= 1 {
			capacity <<= 1
		}
	}
	qu.q = make(chan interface{}, capacity)
	qu.Empty = make(chan struct{})
	qu.out = make(chan interface{}, 1)
	return qu
}

// Capacity() returns the ring buffer capacity

func (qu QueueChan) Capacity() int {
	return cap(qu.q)
}

// length() returns the actual number of elem in the ring buffer

func (qu QueueChan) Length() int {
	return len(qu.q)
}

// Dynamic() shrink the ring buffer capacity on lesser needs

func (qu *QueueChan) Dynamic() {
	qu.dynamic = true
}

// doubleCap() doubles the ring buffer capacity

func (qu *QueueChan) doubleCap() {
	capacity := cap(qu.q) << 1
	//	fmt.Println(capacity)
	close(qu.q)
	tmp := make(chan interface{}, capacity)
	for e := range qu.q {
		tmp <- e
	}
	qu.q = tmp
}

// halfCap() reduce the ring buffer capacity by 50%
// only if qu.dynamic == true (set by Dynamic())

func (qu *QueueChan) halfCap() {
	capacity := cap(qu.q) >> 1
	//	fmt.Println(capacity)
	close(qu.q)
	tmp := make(chan interface{}, capacity)
	for e := range qu.q {
		tmp <- e
	}
	qu.q = tmp
}

// Push(elems ...interface{}) adds the elems
// one by one in the given order to the queues end
// returns Error if no element is provided

func (qu *QueueChan) Push(elems ...interface{}) error {
	if len(elems) == 0 {
		return errors.New("Error: No elem for queue provided.")
	}

	for len(qu.q)+len(elems) > cap(qu.q) {
		qu.doubleCap()
	}
	for _, elem := range elems {
		qu.q <- elem
	}
	return nil
}

// PushTS(elems ...interface{}) threadsafe - adds the elems
// one by one in the given order to the queues end
// returns Error if no element is provided.
// Using push concurently (1) does not guaranty for order of pushes.
// Instead the order of elems pushed together (2)is preserved even if push is used in coroutines.
// But have in mind that PopPushTS used concurrently might break this order, because ring elems will pop and push one by one.
// (1) no order of pushes
// go func(){push("A"))()
// go func(){push("B"))()
// -> can result in any order on the ring buffer:
// .. "A" "B" ..   OR   .. "B" "A" ..
// (2) but ordered within pushes
// go func(){push("A",1,2,3))()
// go func(){push("B",5,6,7))()
// -> .. "B" 5 6 7 .. "A" 1 2 3 ..    OR     .. "A" 1 2 3 .. "B" 5 6 7 ..

func (qu *QueueChan) PushTS(elems ...interface{}) error {
	if len(elems) == 0 {
		return errors.New("Error: No elem for queue provided.")
	}

	qu.Lock()
	defer qu.Unlock()

	for len(qu.q)+len(elems) > cap(qu.q) {
		qu.doubleCap()
	}
	for _, elem := range elems {
		qu.q <- elem
	}
	return nil
}

// Pop() returns and deletes the front elem from the queue

func (qu *QueueChan) Pop() interface{} {
	if len(qu.q) == 0 {
		close(qu.Empty)
		close(qu.out)
		close(qu.q)
		return nil
	}

	e := <-qu.q
	if qu.dynamic && len(qu.q) == cap(qu.q)>>4 {
		qu.halfCap()
	}
	return e
}

// PopTS() threadsafe - returns and deletes the front elem from the queue

func (qu *QueueChan) PopTS() interface{} {
	qu.Lock()
	defer qu.Unlock()

	return qu.Pop()
}

// PopChan() returns a chan with the front elem from the queue and deletes it

func (qu *QueueChan) PopChan() chan interface{} {
	if len(qu.q) == 0 {
		close(qu.Empty)
		close(qu.out)
		close(qu.q)
		return nil
	}

	qu.out <- (<-qu.q)
	if qu.dynamic && len(qu.q) == cap(qu.q)>>4 {
		qu.halfCap()
	}
	return qu.out
}

// PopChanTS() threadsafe - returns a chan with the front elem from the queue and deletes it

func (qu *QueueChan) PopChanTS() chan interface{} {
	qu.Lock()
	defer qu.Unlock()

	return qu.PopChan()
}

// PopPush() returns the front elem from the queue and adds it to the end of the queue (cycling)
func (qu *QueueChan) PopPush() (e interface{}) {
	if len(qu.q) == 0 {
		close(qu.Empty)
		return nil
	}

	e = <-qu.q
	qu.q <- e
	return e
}

// PopPush() threadsafe - returns the front elem from the queue and adds it to the end of the queue (cycling)
func (qu *QueueChan) PopPushTS() (e interface{}) {
	qu.Lock()
	defer qu.Unlock()

	return qu.PopPush()
}

// PopChanPush() returns a chan with the front elem from the queue and adds the elem at the end of the queue (cycling)
func (qu *QueueChan) PopChanPush() chan interface{} {
	if len(qu.q) == 0 {
		close(qu.Empty)
		return nil
	}

	e := <-qu.q
	qu.out <- e
	qu.q <- e
	return qu.out
}

// PopChanPushTS() threadsafe - returns a chan with the front elem from the queue and adds the elem at the end of the queue (cycling)
func (qu *QueueChan) PopChanPushTS() chan interface{} {
	qu.Lock()
	defer qu.Unlock()

	return qu.PopChanPush()
}

// Rotate(nIdx int) rotates by/to nIdx posistion(s)
// Returns nil or an out-of-range error

func (qu *QueueChan) Rotate(np int) error {
	if len(qu.q) == 0 {
		close(qu.Empty)
		return errors.New("Error: can not rotate an empty queue")
	}

	if np < 0 {
		np = np % len(qu.q)
		np = (len(qu.q) + np) % len(qu.q)
	}
	if np > len(qu.q) {
		np = np % len(qu.q)
	}

	for i := 0; i < np; i++ {
		qu.q <- (<-qu.q)
	}
	return nil
}

// RotateTS(nIdx int) threadsafe - rotates by/to nIdx posistion(s)
// Returns nil or an out-of-range error

func (qu *QueueChan) RotateTS(np int) error {
	qu.Lock()
	defer qu.Unlock()

	return qu.Rotate(np)
}
