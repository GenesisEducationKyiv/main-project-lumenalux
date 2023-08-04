package port

type Logger interface {
	Info(...interface{})
	Infof(format string, args ...interface{})

	Debug(...interface{})
	Debugf(format string, args ...interface{})

	Error(...interface{})
	Errorf(format string, args ...interface{})
}
