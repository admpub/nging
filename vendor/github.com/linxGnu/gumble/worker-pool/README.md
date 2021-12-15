# worker-pool

High performance, thread-safe, extendable worker pool.

# Usage

```go
package main

import (
   "log"
   "runtime"
   
   workerpool "github.com/linxGnu/gumble/worker-pool"
)

// task closure
func moduloTask(ctx context.Context, a, b, N uint) *workerpool.Task {
	return workerpool.NewTask(ctx, func(ctx context.Context) (interface{}, error) {
		return modulo(a, b, N), nil
	})
}

// calculate a^b MODULO N
func modulo(a, b uint, N uint) uint {
	switch b {
	case 0:
		return 1 % N
	case 1:
		return a % N
	default:
		if b&1 == 0 {
			t := modulo(a, b>>1, N)
			return uint(uint64(t) * uint64(t) % uint64(N))
		} else {
			t := modulo(a, b>>1, N)
			t = uint(uint64(t) * uint64(t) % uint64(N))
			return uint(uint64(a) * uint64(t) % uint64(N))
		}
	}
}

func main() {
    pool := workerpool.NewPool(nil, workerpool.Option{NumberWorker: runtime.NumCPU()})
	pool.Start()

	// Calculate (1^1 + 2^2 + 3^3 + ... + 1000000^1000000) modulo 1234567
	tasks := make([]*workerpool.Task, 0, 1000000)
	for i := 1; i <= 1000000; i++ {
		task := moduloTask(context.Background(), uint(i), uint(i), 1234567)
		pool.Do(task)
		tasks = append(tasks, task)
	}

	// collect task results
	var s1, s2 uint
	for i := range tasks {
		if result := <-tasks[i].Result(); result.Err != nil {
			log.Fatal(result.Err)
		} else {
			s1 = uint((uint64(s1) + uint64(result.Result.(uint))) % 1234567)
		}
	}

	// sequential computation
	for i := 1; i <= 1000000; i++ {
		s2 = uint((uint64(s2) + uint64(modulo(uint(i), uint(i), 1234567))) % 1234567)
	}
	if s1 != s2 {
		log.Fatal(s1, s2)
    }
    
	pool.Stop()    
}
```
