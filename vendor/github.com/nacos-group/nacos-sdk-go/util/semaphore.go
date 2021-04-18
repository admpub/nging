/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

type Semaphore struct {
	bufSize int
	channel chan int8
}

func NewSemaphore(concurrencyNum int) *Semaphore {
	return &Semaphore{channel: make(chan int8, concurrencyNum), bufSize: concurrencyNum}
}

func (this *Semaphore) TryAcquire() bool {
	select {
	case this.channel <- int8(0):
		return true
	default:
		return false
	}
}

func (this *Semaphore) Acquire() {
	this.channel <- int8(0)
}

func (this *Semaphore) Release() {
	<-this.channel
}

func (this *Semaphore) AvailablePermits() int {
	return this.bufSize - len(this.channel)
}
