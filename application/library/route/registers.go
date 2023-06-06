package route

import "github.com/webx-top/echo"

type Registers []func(g echo.RouteRegister)

func (r *Registers) Register(rg ...func(g echo.RouteRegister)) {
	*r = append(*r, rg...)
}

func (r Registers) Apply(g echo.RouteRegister) {
	for _, rg := range r {
		rg(g)
	}
}
