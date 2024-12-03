package logging

import "context"

type Logger interface {
	Debug(c context.Context, msg string, keysAndValues ...interface{})
	Info(c context.Context, msg string, keysAndValues ...interface{})
	Warn(c context.Context, msg string, keysAndValues ...interface{})
	Error(c context.Context, msg string, keysAndValues ...interface{})
	Fatal(c context.Context, msg string, keysAndValues ...interface{})
}
