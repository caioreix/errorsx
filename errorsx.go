package errorsx

import (
	"fmt"
	"runtime"
	"slices"
	"strings"
)

type ErrorX interface {
	error
	fmt.Stringer

	Caller() string
	Stack() Stack
	Response(fields ...string) map[string]any

	unwrap() ErrorX
	fields() map[string]any
}

type errorX struct {
	stack   Stack
	caller  string
	err     error
	message string
}

var _ ErrorX = (*errorX)(nil)

func (e *errorX) Error() string {
	return stringify(e)
}

func (e *errorX) String() string {
	msg := e.message
	if e.err != nil {
		msg = msg + ": " + e.err.Error()
	}

	return msg
}

func (e *errorX) Response(fields ...string) map[string]any {
	return mapify(e, fields)
}

func (e *errorX) Caller() string {
	return e.caller
}

func (e *errorX) Stack() Stack {
	return e.stack
}

func (e *errorX) unwrap() ErrorX {
	return nil
}

func (e *errorX) fields() map[string]any {
	fields := map[string]any{"message": e.message}
	if e.err != nil {
		fields["error"] = e.err.Error()
	}

	return fields
}

func New(message string) ErrorX {
	return newf(nil, "%s", message)
}

func Newf(format string, args ...any) ErrorX {
	return newf(nil, format, args...)
}

func NewWithError(err error, message string) ErrorX {
	return newf(err, "%s", message)
}

func NewWithErrorf(err error, format string, args ...any) ErrorX {
	return newf(err, format, args...)
}

func newf(err error, format string, args ...any) ErrorX {
	e := &errorX{
		err:     err,
		message: fmt.Sprintf(format, args...),
		caller:  getCaller(2),
		stack:   getStack(4),
	}

	return e
}

func stringify(e ErrorX) string {
	ex := ErrorX(e)
	msgs := make([]string, 1)
	for {
		if eu := ex.unwrap(); eu != nil {
			msgs = append(msgs, ex.String())
			ex = eu
			continue
		}

		msg := ex.String()

		msgs[0] = msg

		msg = strings.Join(msgs, ": ")
		msg = msg + " [" + ex.Caller() + "]"

		return msg
	}
}

func mapify(e ErrorX, fields []string) map[string]any {
	ex := ErrorX(e)
	f := make(map[string]any)
	for {
		if eu := ex.unwrap(); eu != nil {
			mapCopy(f, ex.fields(), fields)
			ex = eu
			continue
		}

		mapCopy(f, ex.fields(), fields)
		if len(fields) == 0 || slices.Contains(fields, "caller") {
			f["caller"] = ex.Caller()
		}

		if len(fields) == 0 || slices.Contains(fields, "caller") {
			f["stack"] = ex.Stack()
		}

		return f
	}
}

func mapCopy(dst, src map[string]any, keys []string) {
	for k, v := range src {
		if len(keys) == 0 || slices.Contains(keys, k) {
			dst[k] = v
		}
	}
}

func getCaller(skip int) string {
	pc, file, line, _ := runtime.Caller(1 + skip)
	return fmt.Sprintf("%s %s:%d", runtime.FuncForPC(pc).Name(), file, line)
}
