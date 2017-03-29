package main

import (
	"fmt"
	"time"

	. "github.com/webx-top/echo/middleware/render/pongo2"
)

func main() {
	t := New(`./template/`)
	t.Init(true)
	//t.SetDebug(true)
	for i := 0; i < 5; i++ {
		ts := time.Now()
		fmt.Printf("==========%v: %v========\\\n", i, ts)
		str := t.Fetch("test", map[string]interface{}{
			`name`: `webx`,
			"test": "times---" + fmt.Sprintf("%v", i),
			"r":    []string{"one", "two", "three"},
		}, nil)
		fmt.Printf("%v\n", str)
		fmt.Printf("==========cost: %v========/\n", time.Now().Sub(ts).Seconds())
	}
}
