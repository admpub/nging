package main

import (
	"fmt"
	"time"

	"github.com/webx-top/echo/middleware/render"
)

type Nested struct {
	Name     string
	Email    string
	Id       int
	HasChild bool
	Children []*Nested
}

func main() {
	tpl := render.New("standard", "./template/")
	tpl.Init(true)
	for i := 0; i < 5; i++ {
		ts := time.Now()
		fmt.Printf("==========%v: %v========\\\n", i, ts)
		str := tpl.Fetch("test", map[string]interface{}{
			"test": "one---" + fmt.Sprintf("%v", i),
			"r":    []string{"one", "two", "three"},
			"nested": []*Nested{
				&Nested{
					Name:     `AAA`,
					Email:    `AAA@webx.top`,
					Id:       1,
					HasChild: true,
					Children: []*Nested{
						&Nested{
							Name:     `AAA1`,
							Email:    `AAA1@webx.top`,
							Id:       11,
							HasChild: true,
							Children: []*Nested{
								&Nested{
									Name:     `AAA11`,
									Email:    `AAA11@webx.top`,
									Id:       111,
									HasChild: false,
								},
							},
						},
					},
				},
				&Nested{
					Name:     `BBB`,
					Email:    `BBB@webx.top`,
					Id:       2,
					HasChild: true,
					Children: []*Nested{
						&Nested{
							Name:     `BBB1`,
							Email:    `BBB1@webx.top`,
							Id:       21,
							HasChild: true,
							Children: []*Nested{
								&Nested{
									Name:     `BBB11`,
									Email:    `BBB11@webx.top`,
									Id:       211,
									HasChild: false,
								},
							},
						},
					},
				},
			},
		}, nil)
		fmt.Printf("%v\n", str)
		fmt.Printf("==========cost: %v========/\n", time.Now().Sub(ts).Seconds())
	}
}
