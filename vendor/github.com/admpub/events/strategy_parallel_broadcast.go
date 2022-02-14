package events

import "golang.org/x/sync/errgroup"

// ParallelBroadcast calls event handlers in separate goroutines
func ParallelBroadcast(event Event, handlers map[Listener]struct{}) error {
	for handler := range handlers {
		go handler.Handle(event)
	}
	return nil
}

// ParallelBroadcastWithReturning calls event handlers in separate goroutines
func ParallelBroadcastWithReturning(event Event, handlers map[Listener]struct{}) error {
	var eg errgroup.Group
	for handler := range handlers {
		h := handler.Handle
		eg.Go(func() error {
			return h(event)
		})
	}
	return eg.Wait()
}
