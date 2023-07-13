//go:build linux

package uring

import (
	"errors"
	"fmt"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

type sq struct {
	buff         []byte
	sqeBuff      []byte
	ringSize     uint64
	kHead        *uint32
	kTail        *uint32
	kRingMask    *uint32
	kRingEntries *uint32
	kFlags       *uint32
	kDropped     *uint32
	kArray       *uint32

	sqeTail, sqeHead uint32
}

func (s *sq) cqNeedFlush() bool {
	return ReadOnceUint32(s.kFlags)&sqCQOverflow != 0
}

type cq struct {
	buff         []byte
	ringSize     uint64
	kHead        *uint32
	kTail        *uint32
	kRingMask    *uint32
	kRingEntries *uint32
	kOverflow    *uint32
	kFlags       uintptr
	cqeBuff      *CQEvent
}

func (c *cq) readyCount() uint32 {
	return SmpLoadAcquireUint32(c.kTail) - SmpLoadAcquireUint32(c.kHead)
}

const MaxEntries uint32 = 1 << 15

//Ring io_uring instance.
type Ring struct {
	fd int

	Params *ringParams

	cqRing *cq
	sqRing *sq
}

var ErrRingSetup = errors.New("ring setup")

type SetupOption func(params *ringParams)

//WithCQSize set CQ max entries count.
func WithCQSize(sz uint32) SetupOption {
	return func(params *ringParams) {
		params.flags = params.flags | setupCQSize
		params.cqEntries = sz
	}
}

//WithIOPoll enable IOPOLL option.
func WithIOPoll() SetupOption {
	return func(params *ringParams) {
		params.flags = params.flags | setupIOPoll
	}
}

//WithAttachedWQ use worker pool from another io_uring instance.
func WithAttachedWQ(fd int) SetupOption {
	return func(params *ringParams) {
		params.flags = params.flags | setupAttachWQ
		params.wqFD = uint32(fd)
	}
}

//WithSQPoll add IORING_SETUP_SQPOLL flag.
//Note, that process must started with root privileges
//or the user should have the CAP_SYS_NICE capability (for kernel version >= 5.11).
func WithSQPoll(threadIdle time.Duration) SetupOption {
	return func(params *ringParams) {
		params.flags = params.flags | setupSQPoll
		params.sqThreadIdle = uint32(threadIdle.Milliseconds())
	}
}

//WithSQThreadCPU bound poll thread to the cpu.
func WithSQThreadCPU(cpu uint32) SetupOption {
	return func(params *ringParams) {
		params.flags = params.flags | setupSQAff
		params.sqThreadCpu = cpu
	}
}

//New create new io_uring instance with. Entries - size of SQ and CQ buffers.
func New(entries uint32, opts ...SetupOption) (*Ring, error) {
	if entries > MaxEntries {
		return nil, fmt.Errorf("%w, entries > MaxEntries", ErrRingSetup)
	}

	params := ringParams{}

	for _, opt := range opts {
		opt(&params)
	}

	fd, err := sysSetup(entries, &params)
	if err != nil {
		return nil, err
	}

	r := &Ring{Params: &params, fd: fd, sqRing: &sq{}, cqRing: &cq{}}
	err = r.allocRing(&params)

	return r, err
}

type Defer func() error

//CreateMany create multiple io_uring instances. Entries - size of SQ and CQ buffers.
//count - the number of io_uring instances. wpCount - the number of worker pools, this value must be a multiple of the entries.
//If workerCount < count worker pool will be shared with setupAttachWQ flag.
func CreateMany(count int, entries uint32, wpCount int, opts ...SetupOption) ([]*Ring, Defer, error) {
	if wpCount > count {
		return nil, nil, errors.New("number of io_uring instances must be greater or equal number of worker pools")
	}
	if count%wpCount != 0 {
		return nil, nil, errors.New("number of worker pools must be a multiple of the entries")
	}

	instancePerPool := count / wpCount

	var defers = map[int]func() error{}
	var rings []*Ring
	for i := 0; i < count; i++ {
		mainRing, err := New(entries, opts...)
		if err != nil {
			for _, closeFn := range defers {
				_ = closeFn()
			}
			return nil, nil, err
		}
		defers[mainRing.fd] = mainRing.Close
		rings = append(rings, mainRing)

		if instancePerPool > 1 {
			for j := 1; j < instancePerPool; j++ {
				additionalRing, err := New(entries, append(opts, WithAttachedWQ(mainRing.fd))...)
				if err != nil {
					for _, closeFn := range defers {
						_ = closeFn()
					}
					return nil, nil, err
				}
				defers[additionalRing.fd] = additionalRing.Close
				rings = append(rings, additionalRing)

				i++
			}
		}
	}

	return rings, func() (err error) {
		for fd, closeFn := range defers {
			cErr := closeFn()
			if cErr != nil {
				err = fmt.Errorf("%w, close ring: %d error: %s", err, fd, cErr.Error())
			}
		}
		return err
	}, nil
}

//Fd io_uring file descriptor.
func (r *Ring) Fd() int {
	return r.fd
}

func (r *Ring) Close() error {
	err := r.freeRing()
	return joinErr(err, syscall.Close(r.fd))
}

var ErrSQOverflow = errors.New("sq ring overflow")

//NextSQE return pointer of the next available SQE in SQ queue.
func (r *Ring) NextSQE() (entry *SQEntry, err error) {
	head := SmpLoadAcquireUint32(r.sqRing.kHead)
	next := r.sqRing.sqeTail + 1

	if next-head <= *r.sqRing.kRingEntries {
		idx := r.sqRing.sqeTail & *r.sqRing.kRingMask * uint32(unsafe.Sizeof(SQEntry{}))
		entry = (*SQEntry)(unsafe.Pointer(&r.sqRing.sqeBuff[idx]))
		r.sqRing.sqeTail = next
	} else {
		err = ErrSQOverflow
	}

	return entry, err
}

type Operation interface {
	PrepSQE(*SQEntry)
	Code() OpCode
}

//QueueSQE adds an operation to the queue SQ.
func (r *Ring) QueueSQE(op Operation, flags uint8, userData uint64) error {
	sqe, err := r.NextSQE()
	if err != nil {
		return err
	}

	op.PrepSQE(sqe)
	sqe.Flags = flags
	sqe.setUserData(userData)
	return nil
}

//Submit new SQEs. Return count of submitted SQEs.
func (r *Ring) Submit() (uint, error) {
	flushed := r.flushSQ()

	var flags uint32

	if !r.needsEnter(&flags) {
		return uint(flushed), nil
	}

	if r.Params.flags&setupIOPoll == 1 {
		flags |= sysRingEnterGetEvents
	}

	consumed, err := sysEnter(r.fd, flushed, 0, flags, nil, true)
	return consumed, err
}

func (r *Ring) needsEnter(flags *uint32) bool {
	if r.Params.flags&setupSQPoll == 0 {
		return true
	}
	if ReadOnceUint32(r.sqRing.kFlags)&sqNeedWakeup != 0 {
		*flags |= sysRingEnterSQWakeup
		return true
	}
	return false
}

var _sizeOfUint32 = unsafe.Sizeof(uint32(0))

func (r *Ring) flushSQ() uint32 {
	mask := *r.sqRing.kRingMask
	tail := SmpLoadAcquireUint32(r.sqRing.kTail)
	subCnt := r.sqRing.sqeTail - r.sqRing.sqeHead

	if subCnt == 0 {
		return tail - SmpLoadAcquireUint32(r.sqRing.kHead)
	}

	for i := subCnt; i > 0; i-- {
		*(*uint32)(unsafe.Add(unsafe.Pointer(r.sqRing.kArray), tail&mask*uint32(_sizeOfUint32))) = r.sqRing.sqeHead & mask
		tail++
		r.sqRing.sqeHead++
	}

	SmpStoreReleaseUint32(r.sqRing.kTail, tail)

	return tail - SmpLoadAcquireUint32(r.sqRing.kHead)
}

type getParams struct {
	submit, waitNr uint32
	flags          uint32
	arg            unsafe.Pointer
	sz             int
}

func (r *Ring) getCQEvents(params getParams) (cqe *CQEvent, err error) {
	for {
		var needEnter = false
		var cqOverflowFlush = false
		var flags uint32
		var available uint32

		available, cqe, err = r.peekCQEvent()
		if err != nil {
			break
		}

		if cqe == nil && params.waitNr == 0 && params.submit == 0 {
			if !r.sqRing.cqNeedFlush() {
				err = syscall.EAGAIN
				break
			}
			cqOverflowFlush = true
		}

		if params.waitNr > available || cqOverflowFlush {
			flags = sysRingEnterGetEvents | params.flags
			needEnter = true
		}

		if params.submit != 0 {
			r.needsEnter(&flags)
			needEnter = true
		}

		if !needEnter {
			break
		}

		var consumed uint
		consumed, err = sysEnter2(r.fd, params.submit, params.waitNr, flags, params.arg, params.sz, false)

		if err != nil {
			break
		}
		params.submit -= uint32(consumed)
		if cqe != nil {
			break
		}
	}

	return cqe, err
}

//WaitCQEventsWithTimeout wait cnt CQEs in CQ. Timeout will be exceeded if no new CQEs in queue.
func (r *Ring) WaitCQEventsWithTimeout(cnt uint32, timeout time.Duration) (cqe *CQEvent, err error) {
	if r.Params.ExtArgFeature() {
		ts := syscall.NsecToTimespec(timeout.Nanoseconds())
		arg := newGetEventsArg(uintptr(unsafe.Pointer(nil)), numSig/8, uintptr(unsafe.Pointer(&ts)))

		cqe, err = r.getCQEvents(getParams{
			submit: 0,
			waitNr: cnt,
			flags:  sysRingEnterExtArg,
			arg:    unsafe.Pointer(arg),
			sz:     int(unsafe.Sizeof(getEventsArg{})),
		})

		runtime.KeepAlive(arg)
		runtime.KeepAlive(ts)
		return cqe, err
	}

	var toSubmit uint32

	var sqe *SQEntry
	sqe, err = r.NextSQE()
	if err != nil {
		_, err = r.Submit()
		if err != nil {
			return nil, err
		}

		sqe, err = r.NextSQE()
		if err != nil {
			return nil, err
		}
	}

	op := Timeout(timeout)
	op.PrepSQE(sqe)
	sqe.setUserData(libUserDataTimeout)
	toSubmit = r.flushSQ()

	return r.getCQEvents(getParams{
		submit: toSubmit,
		waitNr: cnt,
		arg:    unsafe.Pointer(nil),
		sz:     numSig / 8,
	})
}

//WaitCQEvents wait cnt CQEs in CQ.
func (r *Ring) WaitCQEvents(cnt uint32) (cqe *CQEvent, err error) {
	return r.getCQEvents(getParams{
		submit: 0,
		waitNr: cnt,
		arg:    unsafe.Pointer(nil),
		sz:     numSig / 8,
	})
}

//SubmitAndWaitCQEvents submit new SQEs in SQE and wait cnt CQEs in CQ. Return first available CQE.
func (r *Ring) SubmitAndWaitCQEvents(cnt uint32) (cqe *CQEvent, err error) {
	return r.getCQEvents(getParams{
		submit: r.flushSQ(),
		waitNr: cnt,
		arg:    unsafe.Pointer(nil),
		sz:     numSig / 8,
	})
}

//PeekCQE return first available CQE from CQ.
func (r *Ring) PeekCQE() (*CQEvent, error) {
	return r.WaitCQEvents(0)
}

//SeenCQE dequeue CQ on 1 entry.
func (r *Ring) SeenCQE(cqe *CQEvent) {
	r.AdvanceCQ(1)
}

//AdvanceCQ dequeue CQ on n entry.
func (r *Ring) AdvanceCQ(n uint32) {
	SmpStoreReleaseUint32(r.cqRing.kHead, *r.cqRing.kHead+n)
}

func (r *Ring) peekCQEvent() (uint32, *CQEvent, error) {
	mask := *r.cqRing.kRingMask
	var cqe *CQEvent
	var available uint32

	var err error
	for {
		tail := SmpLoadAcquireUint32(r.cqRing.kTail)
		head := SmpLoadAcquireUint32(r.cqRing.kHead)

		cqe = nil
		available = tail - head
		if available == 0 {
			break
		}

		cqe = (*CQEvent)(unsafe.Add(unsafe.Pointer(r.cqRing.cqeBuff), uintptr(head&mask)*unsafe.Sizeof(CQEvent{})))

		if !r.Params.ExtArgFeature() && cqe.UserData == libUserDataTimeout {
			if cqe.Res < 0 {
				err = cqe.Error()
			}
			r.SeenCQE(cqe)
			if err == nil {
				continue
			}
			cqe = nil
		}
		break
	}

	return available, cqe, err
}

func (r *Ring) peekCQEventBatch(buff []*CQEvent) int {
	ready := r.cqRing.readyCount()
	count := min(uint32(len(buff)), ready)

	if ready != 0 {
		head := SmpLoadAcquireUint32(r.cqRing.kHead)
		mask := SmpLoadAcquireUint32(r.cqRing.kRingMask)

		last := head + count
		for i := 0; head != last; head, i = head+1, i+1 {
			buff[i] = (*CQEvent)(unsafe.Add(unsafe.Pointer(r.cqRing.cqeBuff), uintptr(head&mask)*unsafe.Sizeof(CQEvent{})))
		}
	}
	return int(count)
}

func min(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}

//PeekCQEventBatch fill buffer by available CQEs. Return count of filled CQEs.
func (r *Ring) PeekCQEventBatch(buff []*CQEvent) int {
	n := r.peekCQEventBatch(buff)
	if n == 0 {
		if r.sqRing.cqNeedFlush() {
			_, _ = sysEnter(r.fd, 0, 0, sysRingEnterGetEvents, nil, false)
			n = r.peekCQEventBatch(buff)
		}
	}

	return n
}

func joinErr(err1, err2 error) error {
	if err1 == nil {
		return err2
	}
	if err2 == nil {
		return err1
	}

	return fmt.Errorf("multiple errors: %w and %s", err1, err2.Error())
}
