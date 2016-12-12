/*

   Copyright 2016 Wenhui Shen <www.webx.top>

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package events

import (
	"sync"
)

var Events *EventsInstance = NewEvents()

func NewEvents() *EventsInstance {
	return &EventsInstance{
		listeners: make(map[string][]func(func(bool), ...interface{})),
		lock:      new(sync.RWMutex),
	}
}

//并发执行事件
func GoEvent(eventName string, next func(bool), sessions ...interface{}) {
	Events.GoExecute(eventName, next, sessions...)
}

//顺序执行事件
func Event(eventName string, next func(bool), sessions ...interface{}) {
	Events.Execute(eventName, next, sessions...)
}

//删除事件
func DelEvent(eventName string) {
	Events.Delete(eventName)
}

//添加事件
func AddEvent(eventName string, handler func(func(bool), ...interface{})) {
	Events.Register(eventName, handler)
}

type EventsInstance struct {
	listeners map[string][]func(func(bool), ...interface{})
	lock      *sync.RWMutex
}

/*
注册事件
[Examle:]
Events.Register("AfterResponse", func(next func(bool),session ...interface{}) {
	log.Println("Got AfterResponse event!")
	isSuccess := true
	next(isSuccess) //这里的next函数无论什么情况下必须执行。
})

采用不同的方式执行事件时，此处的next函数的作用也是不同的：
1、在并发执行事件的时候，next函数的作用是通知程序我已经执行完了(不理会这一步是否执行成功)；
2、在顺序执行事件的时候，next函数的作用是通知程序是否继续执行下一步，next(true)是继续执行下一步，next(false)是终止执行下一步
*/
func (e *EventsInstance) Register(eventName string, handler func(func(bool), ...interface{})) {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.listeners == nil {
		e.listeners = make(map[string][]func(func(bool), ...interface{}))
	}
	_, ok := e.listeners[eventName]
	if !ok {
		e.listeners[eventName] = make([]func(func(bool), ...interface{}), 0)
	}
	e.listeners[eventName] = append(e.listeners[eventName], handler)
}

func (e *EventsInstance) Delete(eventName string) {
	e.lock.Lock()
	defer e.lock.Unlock()
	if e.listeners == nil {
		return
	}
	_, ok := e.listeners[eventName]
	if ok {
		delete(e.listeners, eventName)
	}
}

/*
并发执行事件
[Examle 1:]
Events.GoExecute("AfterHandler", func(_ bool) {//此匿名函数在本事件的最后执行
	session.Response.Send()
	session.Response.Close()
}, session)

[Examle 2:]
Events.Execute("AfterResponse", func(_ bool) {}, session)
*/
func (e *EventsInstance) GoExecute(eventName string, next func(bool), sessions ...interface{}) {
	if e.listeners == nil {
		next(true)
		return
	}
	c := make(chan int)
	n := 0
	e.lock.RLock()
	defer e.lock.RUnlock()
	if l, ok := e.listeners[eventName]; ok {
		if len(l) > 0 {
			for _, h := range l {
				n++
				//h 的原型为 func(interface{}, func(bool))
				go h(func(_ bool) {
					c <- 1
				}, sessions...)
			}
		}
	}
	for n > 0 {
		i := <-c
		if i == 1 {
			n--
		}
	}
	if next == nil {
		return
	}
	next(true)
}

/**
 * 顺序执行事件
 */
func (e *EventsInstance) Execute(eventName string, next func(bool), sessions ...interface{}) {
	if e.listeners == nil {
		next(true)
		return
	}
	e.lock.RLock()
	defer e.lock.RUnlock()
	var nextStep bool = false
	if l, ok := e.listeners[eventName]; ok {
		if len(l) > 0 {
			for _, h := range l {
				h(func(ok bool) {
					nextStep = ok
				}, sessions...)
				//一旦传入false，后面的全部忽略执行
				if !nextStep {
					break
				}
			}
		} else {
			nextStep = true
		}
	} else {
		nextStep = true
	}
	if next == nil {
		return
	}
	next(nextStep)
}
