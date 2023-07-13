package reactor

type Logger interface {
	Log(keyvals ...interface{})
}

type nopLogger struct {
}

func (n *nopLogger) Log(keyvals ...interface{}) {
}
