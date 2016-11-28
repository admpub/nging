#queueChan

Go/Golang fast FIFO queue or ring buffer based on channels.


### Summary
Fast FIFO queue (*Push()*-*Pop()*) using a channel for underlying structure. *PopChan()* provides a channel with the element instead of the element itself for use within select. 

Used as a fast ring-buffer too by adding the extracted front elements back at the queues end: *PopPush()* or *PopChanPush()*

Ring-buffer head might be shifted up or down by n steps: *Rotate(n)*

Threadsafe versions of functions for concurrent use (in goroutines): *PushTS()*, *PopTS()*, *PopChanTS()*, *PopPushTS()*, *PopChanPushTS()*, *RotateTS()*

Have in mind that QueueChan will close down channels if the buffer is empty allowing for GC and graceful exit of consuming code - This needs to be considered in the porgram logic. I would be happy to get suggestions on this topic.

Threadsafe calls are 6-7 times slower (because of the blocking Mutex on the struct to prevent race conditions) and by nature don't quaranty for preserved serial order of pushs or pops. The design is streamlined and fast - concurrent use of non threadsafe method calls are very likely to cause data losses and panic. 

Anyway elements pushed together in one call *Push(a,b,c)* will preserve order a - b - c in the buffer even with concurrent use of *PushTS(a,b,c)*. 

In the end you might come up with mixing threadsafe and non-threadsafe calls depending on different phases of a problem solution, i.e. nonconcurrent *Push()* (to produce an ordered element pile), nonconcurrent *PopPush()* for an ordered logging and lastly a concurrent *PopTS()* to work them off regardless of order in a later step.

### Usage
Get it as usual with 

```go 
	go get github.com/AndreasBriese/queueChan
``` 

and import with 

```go 
	import (
		...
		"github.com/AndreasBriese/queueChan"
	)
``` 

Instanziate by your taste and with/without predefined length (default is capacity=16; will always use the next higher pow2 (1000->1024)).

```go
	// Declaration
	var qc1 queueChan.QueueChan = queueChan.New(1000)

	qc2 := queueChan.New(1000)
	qc3 := (&queueChan.QueueChan{}).New()
	
	fmt.Printf("number of elements:%v  capacity: %v  default capacity: %v", qc1.Length(), qc2.Capacity(), qc3.Capacity())
	
	qc3.Dynamic() // auto shrink from now on, if element are taken away
``` 

QueueChan will allways scale up automatically (doubling) if the capacity is exhausted, but scaling costs process time and memory while processed (the bigger the pile of elements the more, because existing elements need to be copied over to a new temporary channel) and using a predefined approximated capacity will leverage this. 

On default the capacity will not shrink when elements are taken from the pile since this is expemsive too. Use method *Dynamic()* to enable continuous shrinking on your queue, if you want to free memory depending on workload at some stage of processing.    

See xxx\_test.go for more use information.

(c) 2015 Andreas Briese,  eduToolbox@BriC GmbH, Sarstedt
MIT License 
