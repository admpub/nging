package main

import (
	"fmt"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/handler/mvc"
)

type Index struct {
	ping    mvc.Mapper `webx:"ping"`
	noafter mvc.Mapper
	echo.Context
	exit bool
}

func (a *Index) Init(ctx echo.Context) error {
	a.Context = ctx
	return nil
}

func (a *Index) Before() error {
	fmt.Println(`-------------->Before`)
	return nil
}

func (a *Index) Ping() error {
	fmt.Println(`-------------->Ping`)
	return a.String(`pong`)
}

func (a *Index) Noafter() error {
	fmt.Println(`-------------->Noafter`)
	a.exit = true
	return a.String(`pong`)
}

func (a *Index) After() error {
	fmt.Println(`-------------->After`)
	return nil
}

func (a *Index) IsExit() bool {
	return a.exit
}

func main() {
	s := mvc.New(`test`)
	m := s.Module()
	m.Register(`/`, func(ctx echo.Context) error {
		return ctx.String(`Hello world.`)
	})
	m.Use(&Index{})
	s.Run(`:8181`)
}
