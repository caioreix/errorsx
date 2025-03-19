package errorsx

import (
	"fmt"
	"runtime"
	"strings"
)

type ErrorX interface {
	error
	fmt.Stringer

	Caller() string
	unwrap() ErrorX
}

type errorX struct {
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

func (e *errorX) Caller() string {
	return e.caller
}

func (e *errorX) unwrap() ErrorX {
	return nil
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

func getCaller(skip int) string {
	pc, _, line, _ := runtime.Caller(1 + skip)
	return fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), line)
}
