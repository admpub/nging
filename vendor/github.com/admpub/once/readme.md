## sync-once

sync-once is similar to [sync.Once](https://golang.org/pkg/sync/#Once) of the standard library.

It also has struct `Once` but has two additional methods `DoForce()` and `Reset()`.

### why

While writing a web application I needed to reload configurations that were calculated using sync.Once. 

But `sync.Once` provide a single method `Do()` that executes the function only once.

To get around this I wrote this package.

### usage

```
import (
    sync "github.com/admpub/once"
)

func main() {
    o := new(sync.Once)
    
    // This will work similar to the Once.Do(f) method of the sync package. The function f is only called once
    o.Do(loadConfig())

    // To call the function either for the first time or again you need to use the sync.DoForce() method
    // This will execute irrespective of weather o.Do() was called earlier or not and mark o (Once) as done.
    // Call to o.Do() after o.DoForce() will not execute the function.
    o.DoForce(loadConfig())

    // To reset o (sync.Once) you need to call the Reset() method.
    // This will mark o as not executed but will no call the Do() method. You need to call Do() or DoForce() after this.
    // Calls to Do() and DoForce() after this will work as described above.
    o.Reset()

}

// load config from a static file or any other operation that is usually performed only once
func loadConfig() {
    // Do the work here
}
``` 
