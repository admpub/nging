package internal

import (
	"fmt"
	"time"
)

type myTimer struct {
	start time.Time
	end   time.Time
}

func NewMyTimer() *myTimer {
	return &myTimer{
		start: time.Now(),
	}
}

func (mt *myTimer) Stop() {
	mt.end = time.Now()
}
func (mt *myTimer) UsedSecond() string {
	return fmt.Sprintf("%f s", mt.end.Sub(mt.start).Seconds())
}
