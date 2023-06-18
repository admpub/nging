package gerberos

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"testing"
	"time"
)

var (
	errFault = errors.New("fault")
)

type testFaultyExecutor struct {
	// Actuator
	name string
	args []string

	// Effect
	output   string
	exitCode int
	err      error
}

func (e *testFaultyExecutor) Execute(name string, args ...string) (string, int, error) {
	return e.ExecuteWithStd(nil, nil, name, args...)
}

func (e *testFaultyExecutor) ExecuteWithStd(stdin io.Reader, stdout io.Writer, name string, args ...string) (string, int, error) {
	if name == e.name && reflect.DeepEqual(args, e.args) {
		return e.output, e.exitCode, e.err
	}

	de := &defaultExecutor{}
	return de.ExecuteWithStd(stdin, stdout, name, args...)
}

func NewTestFaultyExecutor(output string, exitCode int, err error, name string, args ...string) *testFaultyExecutor {
	return &testFaultyExecutor{
		name:     name,
		args:     args,
		output:   output,
		exitCode: exitCode,
		err:      err,
	}
}

func testError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Error("expected error")
	}
}

func testNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
}

func newTestConfiguration() (*Configuration, error) {
	c := &Configuration{}

	return c, c.ReadFile("test/configuration.toml")
}

func newTestRunner() (*Runner, error) {
	c, err := newTestConfiguration()
	if err != nil {
		return nil, err
	}

	return NewRunner(c), nil
}

func newTestValidRule() *Rule {
	return &Rule{
		Action:      []string{"ban", "1h"},
		Regexp:      []string{`%ip%\s%id%`},
		Source:      []string{"test"},
		Aggregate:   []string{"1s", `a\s%id%`, `%id%\sb`},
		Occurrences: []string{"5", "10s"},

		name: "test",
	}
}

func newTestOccurrences() *occurrences {
	return newOccurrences(100*time.Millisecond, 10)
}

func testCountChildren() (int, error) {
	cmd := exec.Command("pgrep", "-P", strconv.Itoa(os.Getpid()))
	b, err := cmd.CombinedOutput()
	if err != nil {
		return 0, nil // Probably no children
	}
	child := 0
	for i := range b {
		if b[i] == '\n' {
			child++
		}
	}
	return child, nil
}
