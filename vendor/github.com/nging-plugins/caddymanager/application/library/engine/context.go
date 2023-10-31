package engine

import (
	"context"
	"io"
)

type ContextKey string

const CtxKeyStdout ContextKey = `stdout`
const CtxKeyStderr ContextKey = `stderr`

func WithStdout(c context.Context, w io.Writer) context.Context {
	return context.WithValue(c, CtxKeyStdout, w)
}

func WithStderr(c context.Context, w io.Writer) context.Context {
	return context.WithValue(c, CtxKeyStderr, w)
}

func WithStdoutStderr(c context.Context, w ...io.Writer) context.Context {
	length := len(w)
	if length == 0 {
		return c
	}
	c = WithStdout(c, w[0])
	if length > 1 {
		c = WithStderr(c, w[1])
	}
	return c
}

func GetCtxStdout(c context.Context) io.Writer {
	v, _ := c.Value(CtxKeyStdout).(io.Writer)
	return v
}

func GetCtxStderr(c context.Context) io.Writer {
	v, _ := c.Value(CtxKeyStderr).(io.Writer)
	return v
}
