package verbose

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"sync"
)

type contextualError struct {
	err error
	key interface{}
	val string
}

func (e contextualError) Error() string {
	return e.err.Error()
}

func (e contextualError) Unwrap() error {
	return e.err
}

var traceKey struct{}

type stackTraceLogger struct {
	mux sync.Mutex
	log func(error) string
}

var defaultStackTraceLogger = stackTraceLogger{
	log: func(err error) string {
		if err == nil {
			return ""
		}
		b := bytes.NewBufferString(err.Error())
		b.WriteString(":\n")
		for skip := 2; skip < 7; skip++ {
			pc, file, line, ok := runtime.Caller(skip)
			if !ok {
				break
			}
			fn := runtime.FuncForPC(pc).Name()
			fmt.Fprintf(b, "    %s:%d %s\n", file, line, fn)
		}
		return b.String()
	},
}

func SetDigger(loggerFn func(error) string) {
	if loggerFn == nil {
		return
	}
	defaultStackTraceLogger.mux.Lock()
	defer defaultStackTraceLogger.mux.Unlock()
	defaultStackTraceLogger.log = loggerFn
}

func WithValue(err error, key, val string) error {
	return &contextualError{
		err: err,
		key: key,
		val: val,
	}
}

func Value(err error, key string) (string, bool) {
	return value(err, key)
}

func value(err error, key interface{}) (string, bool) {
	var ctxErr *contextualError
	if !errors.As(err, &ctxErr) {
		return "", false
	}
	if ctxErr.key == key {
		return ctxErr.val, true
	}
	return value(ctxErr.err, key)
}

func WithStackTrace(err error) error {
	return &contextualError{
		err: err,
		key: traceKey,
		val: defaultStackTraceLogger.log(err),
	}
}

func StackTrace(err error) (string, bool) {
	var b bytes.Buffer
	find := false
	for {
		var ctxErr *contextualError
		if !errors.As(err, &ctxErr) {
			break
		}
		if ctxErr.key == traceKey {
			find = true
			b.WriteString(ctxErr.val)
		}
		err = ctxErr.err
	}
	return b.String(), find
}
