package ppnocolor

import (
	"io"

	"github.com/admpub/pp"
)

func New(callerLevel ...int) *pp.PrettyPrinter {
	a := pp.New(callerLevel...)
	a.SetColoringEnabled(false)
	return a
}

var Default = New(3)

// Print prints given arguments.
func Print(a ...interface{}) (n int, err error) {
	return Default.Print(a...)
}

// Printf prints a given format.
func Printf(format string, a ...interface{}) (n int, err error) {
	return Default.Printf(format, a...)
}

// Println prints given arguments with newline.
func Println(a ...interface{}) (n int, err error) {
	return Default.Println(a...)
}

// Sprint formats given arguments and returns the result as string.
func Sprint(a ...interface{}) string {
	return Default.Sprint(a...)
}

// Sprintf formats with pretty print and returns the result as string.
func Sprintf(format string, a ...interface{}) string {
	return Default.Sprintf(format, a...)
}

// Sprintln formats given arguments with newline and returns the result as string.
func Sprintln(a ...interface{}) string {
	return Default.Sprintln(a...)
}

// Fprint prints given arguments to a given writer.
func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	return Default.Fprint(w, a...)
}

// Fprintf prints format to a given writer.
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	return Default.Fprintf(w, format, a...)
}

// Fprintln prints given arguments to a given writer with newline.
func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	return Default.Fprintln(w, a...)
}

// Errorf formats given arguments and returns it as error type.
func Errorf(format string, a ...interface{}) error {
	return Default.Errorf(format, a...)
}

// Fatal prints given arguments and finishes execution with exit status 1.
func Fatal(a ...interface{}) {
	Default.Fatal(a...)
}

// Fatalf prints a given format and finishes execution with exit status 1.
func Fatalf(format string, a ...interface{}) {
	Default.Fatalf(format, a...)
}

// Fatalln prints given arguments with newline and finishes execution with exit status 1.
func Fatalln(a ...interface{}) {
	Default.Fatalln(a...)
}

// Change Print* functions' output to a given writer.
// For example, you can limit output by ENV.
//
//	func init() {
//		if os.Getenv("DEBUG") == "" {
//			pp.SetDefaultOutput(ioutil.Discard)
//		}
//	}
func SetDefaultOutput(o io.Writer) {
	Default.SetOutput(o)
}

// GetOutput returns pp's default output.
func GetDefaultOutput() io.Writer {
	return Default.GetOutput()
}

// Change Print* functions' output to default one.
func ResetDefaultOutput() {
	Default.ResetOutput()
}

// SetColorScheme takes a colorscheme used by all future Print calls.
func SetColorScheme(scheme pp.ColorScheme) {
	Default.SetColorScheme(scheme)
}

// ResetColorScheme resets colorscheme to default.
func ResetColorScheme() {
	Default.ResetColorScheme()
}

// SetMaxDepth sets the printer's Depth, -1 prints all
func SetDefaultMaxDepth(v int) {
	Default.SetDefaultMaxDepth(v)
}
