package flock

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"
)

func TestLock(t *testing.T) {
	testFilePath, _ := os.Getwd()
	wg := sync.WaitGroup{}
	fpEx, err := os.Open(testFilePath)
	if err != nil {
		panic(err)
	}
	defer fpEx.Close()
	err = LockEx(fpEx)
	if err != nil {
		fmt.Println(`lock error: `, err.Error())
		return
	}
	go func() {
		time.Sleep(2 * time.Second)
		Unlock(fpEx)
	}()
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			fp, err := os.Open(testFilePath)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			defer fp.Close()
			err = LockBlock(fp) // LockBlock 需要等待fpEx解锁之后才能继续
			if err != nil {
				fmt.Println(`lock error: `, err.Error())
				return
			}
			fmt.Printf("output : %d\n", num)
		}(i)
	}
	wg.Wait()
	time.Sleep(2 * time.Second)
}
