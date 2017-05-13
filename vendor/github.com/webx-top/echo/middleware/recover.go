package middleware

import (
	"fmt"
	"runtime"

	"github.com/webx-top/echo"
)

type (
	// RecoverConfig defines the config for Recover middleware.
	RecoverConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echo.Skipper

		// Size of the stack to be printed.
		// Optional. Default value 4KB.
		StackSize int `json:"stack_size"`

		// DisableStackAll disables formatting stack traces of all other goroutines
		// into buffer after the trace for the current goroutine.
		// Optional. Default value false.
		DisableStackAll bool `json:"disable_stack_all"`

		// DisablePrintStack disables printing stack trace.
		// Optional. Default value as false.
		DisablePrintStack bool `json:"disable_print_stack"`
	}
)

var (
	// DefaultRecoverConfig is the default Recover middleware config.
	DefaultRecoverConfig = RecoverConfig{
		Skipper:           echo.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
	}
)

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover() echo.Middleware {
	return RecoverWithConfig(DefaultRecoverConfig)
}

// RecoverWithConfig returns a Recover middleware with config.
// See: `Recover()`.
func RecoverWithConfig(config RecoverConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = DefaultRecoverConfig.StackSize
	}

	return func(next echo.Handler) echo.Handler {
		return echo.HandlerFunc(func(c echo.Context) error {
			if config.Skipper(c) {
				return next.Handle(c)
			}

			defer func() {
				if r := recover(); r != nil {
					panicErr := echo.NewPanicError(r, nil)
					panicErr.SetDebug(c.Echo().Debug())
					var err error
					switch r := r.(type) {
					case error:
						err = r
					default:
						err = fmt.Errorf("%v", r)
					}
					if config.DisableStackAll {
						c.Error(panicErr.SetError(err))
						return
					}
					content := "[PANIC RECOVER] " + err.Error()
					for i := 1; len(content) < config.StackSize; i++ {
						pc, file, line, ok := runtime.Caller(i)
						if !ok {
							break
						}
						t := &echo.Trace{
							File: file,
							Line: line,
							Func: runtime.FuncForPC(pc).Name(),
						}
						panicErr.AddTrace(t)
						content += "\n" + fmt.Sprintf(`%v:%v`, file, line)
					}
					panicErr.SetErrorString(content)
					c.Logger().Error(panicErr)
					c.Error(panicErr)
				}
			}()
			return next.Handle(c)
		})
	}
}
