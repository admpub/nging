package echo

type (
	// Skipper defines a function to skip middleware. Returning true skips processing
	// the middleware.
	Skipper func(c Context) bool
)

// defaultSkipper returns false which processes the middleware.
func DefaultSkipper(c Context) bool {
	return false
}
