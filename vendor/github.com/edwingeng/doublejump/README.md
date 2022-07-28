# Overview
Doublejump is a revamped [Google's jump](https://arxiv.org/pdf/1406.2294.pdf) consistent hash. It overcomes the shortcoming of the original design - being unable to remove nodes. Here is [how it works](https://docs.google.com/presentation/d/e/2PACX-1vTHyFGUJ5CBYxZTzToc_VKxP_Za85AeZqQMNGLXFLP1tX0f9IF_z3ys9-pyKf-Jj3iWpm7dUDDaoFyb/pub?start=false&loop=false&delayms=3000).

# Benchmark
```
DoubleJump/10-nodes               49276861       22.3 ns/op        0 B/op      0 allocs/op
DoubleJump/100-nodes              33304191       34.9 ns/op        0 B/op      0 allocs/op
DoubleJump/1000-nodes             25261296       46.3 ns/op        0 B/op      0 allocs/op

StathatConsistent/10-nodes         4780832      273.5 ns/op       80 B/op      2 allocs/op
StathatConsistent/100-nodes        4059537      291.8 ns/op       80 B/op      2 allocs/op
StathatConsistent/1000-nodes       3132294      367.6 ns/op       80 B/op      2 allocs/op

SerialxHashring/10-nodes           2766384      455.7 ns/op      152 B/op      5 allocs/op
SerialxHashring/100-nodes          2500936      487.6 ns/op      152 B/op      5 allocs/op
SerialxHashring/1000-nodes         2254138      560.0 ns/op      152 B/op      5 allocs/op
```

# Installation

### V1
```shell
## If golang version <= 1.17
go get -u github.com/edwingeng/doublejump
```

### V2
```shell
## If golang version >= 1.18
go get -u github.com/edwingeng/doublejump/v2
```

# Examples

### V1
```go
// If golang version <= 1.17
import "github.com/edwingeng/doublejump"

func Example() {
    h := NewHash()
    for i := 0; i < 10; i++ {
        h.Add(fmt.Sprintf("node%d", i))
    }

    fmt.Println(h.Len())
    fmt.Println(h.LooseLen())

    fmt.Println(h.Get(1000))
    fmt.Println(h.Get(2000))
    fmt.Println(h.Get(3000))

    h.Remove("node3")
    fmt.Println(h.Len())
    fmt.Println(h.LooseLen())

    fmt.Println(h.Get(1000))
    fmt.Println(h.Get(2000))
    fmt.Println(h.Get(3000))

    // Output:
    // 10
    // 10
    // node9
    // node2
    // node3
    // 9
    // 10
    // node9
    // node2
    // node0
}
```

### V2
```go
// If golang version >= 1.18
import "github.com/edwingeng/doublejump/v2"

func Example() {
    h := NewHash[string]()
    for i := 0; i < 10; i++ {
        h.Add(fmt.Sprintf("node%d", i))
    }

    fmt.Println(h.Len())
    fmt.Println(h.LooseLen())

    fmt.Println(h.Get(1000))
    fmt.Println(h.Get(2000))
    fmt.Println(h.Get(3000))

    h.Remove("node3")
    fmt.Println(h.Len())
    fmt.Println(h.LooseLen())

    fmt.Println(h.Get(1000))
    fmt.Println(h.Get(2000))
    fmt.Println(h.Get(3000))

    // Output:
    // 10
    // 10
    // node9 true
    // node2 true
    // node3 true
    // 9
    // 10
    // node9 true
    // node2 true
    // node0 true
}
```

# Acknowledgements
The implementation of the original algorithm is credited to [dgryski](https://github.com/dgryski/go-jump).
