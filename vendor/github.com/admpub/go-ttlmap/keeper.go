package ttlmap

import "time"

/*
lock
*/

type keeper struct {
	store        *store
	updating     bool
	drained      bool
	updateChan   chan struct{}
	drainingChan chan struct{}
	drainChan    chan struct{}
	doneChan     chan struct{}
}

func newKeeper(store *store) *keeper {
	return &keeper{
		store:        store,
		updateChan:   make(chan struct{}, 1),
		drainingChan: make(chan struct{}),
		drainChan:    make(chan struct{}, 1),
		doneChan:     make(chan struct{}),
	}
}

func (k *keeper) run() {
	defer close(k.doneChan)
	defer k.drain()
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-k.drainingChan:
			return
		case <-k.updateChan:
			k.update(timer, false)
		case <-timer.C:
			k.update(timer, true)
		}
	}
}

func (k *keeper) signalDrain() {
	select {
	case k.drainChan <- struct{}{}:
		close(k.drainingChan)
	default:
	}
}

func (k *keeper) signalUpdate() {
	if !k.updating {
		k.updating = true
		select {
		case k.updateChan <- struct{}{}:
		default:
		}
	}
}

func (k *keeper) update(timer *time.Timer, evict bool) {
	k.store.Lock()
	if evict {
		k.store.evictExpired()
	}
	k.updating = false
	duration, ok := k.nextTTL()
	k.store.Unlock()
	if ok {
		timer.Reset(duration)
	} else {
		timer.Stop()
	}
}

func (k *keeper) nextTTL() (time.Duration, bool) {
	pqi := k.store.pq.peek()
	if pqi == nil {
		return 0, false
	}
	duration := pqi.item.TTL()
	if duration < 0 {
		duration = 0
	}
	return duration, true
}

func (k *keeper) drain() {
	k.store.Lock()
	k.drained = true
	k.store.drain()
	k.store.Unlock()
}
