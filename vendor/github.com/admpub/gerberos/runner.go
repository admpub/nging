package gerberos

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Runner struct {
	Configuration      *Configuration
	backend            Backend
	respawnWorkerDelay time.Duration
	respawnWorkerChan  chan *Rule
	Executor           Executor
	stop               context.CancelFunc
	stopped            context.Context
}

func (rn *Runner) Initialize() error {
	if rn.Configuration == nil {
		return errors.New("configuration has not been set")
	}

	// Backend
	switch rn.Configuration.Backend {
	case "":
		return errors.New("missing configuration value for backend")
	default:
		bfn, ok := backends[rn.Configuration.Backend]
		if !ok {
			return fmt.Errorf("unknown backend: %s", rn.Configuration.Backend)
		}
		rn.backend = bfn(rn)
	}
	if err := rn.backend.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize backend: %w", err)
	}

	// Rules
	for n, r := range rn.Configuration.Rules {
		r.name = n
		if err := r.initialize(rn); err != nil {
			return fmt.Errorf(`failed to initialize Rule "%s": %s`, n, err)
		}
	}

	return nil
}

func (rn *Runner) Finalize() error {
	if err := rn.backend.Finalize(); err != nil {
		return fmt.Errorf("failed to finalize backend: %w", err)
	}

	return nil
}

func (rn *Runner) spawnWorker(r *Rule, requeue bool) {
	go func() {
		select {
		case <-rn.stopped.Done():
		default:
			r.worker(requeue)
		}
	}()
	log.Printf("%s: spawned worker", r.name)
}

func (rn *Runner) Run(requeueWorkers bool) {
	for _, r := range rn.Configuration.Rules {
		rn.spawnWorker(r, requeueWorkers)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(signalChan)

	go func() {
		for {
			select {
			case r := <-rn.respawnWorkerChan:
				time.Sleep(rn.respawnWorkerDelay)
				rn.spawnWorker(r, requeueWorkers)
			case <-rn.stopped.Done():
				return
			}
		}
	}()

	select {
	case <-rn.stopped.Done():
	case s := <-signalChan:
		log.Printf("received signal: %s", s)
		rn.stop()
	}
}

func (rn *Runner) Stop() {
	rn.stop()
}

func (rn *Runner) Ban(ip string, ipv6 bool, d time.Duration) error {
	return rn.backend.Ban(ip, ipv6, d)
}

func NewRunner(c *Configuration) *Runner {
	ctx, cancel := context.WithCancel(context.Background())
	return &Runner{
		Configuration:      c,
		respawnWorkerDelay: 5 * time.Second,
		respawnWorkerChan:  make(chan *Rule),
		Executor:           &defaultExecutor{},
		stop:               cancel,
		stopped:            ctx,
	}
}
