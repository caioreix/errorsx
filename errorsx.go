package errorsx

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
)

const (
	errorXDefaultType = "default"
	errorXHttpType    = "http"
)

type ErrorX interface {
	error

	Wrap(err error) ErrorX
	Skip(n int) ErrorX
}

var _ ErrorX = (*errorX)(nil)

type errorX struct {
	eType string

	caller string
	skip   int

	err error

	message string

	errorXHttp
}

type errorXHttp struct {
	status int
}

func (e *errorX) Error() string {
	logs := []string{e.message}

	if e.eType == errorXHttpType {
		logs = append(logs, fmt.Sprintf("status: %d", e.status))
	}

	if e.err != nil {
		logs = append(logs, e.err.Error())
	}

	msg := strings.Join(logs, ": ")

	return fmt.Sprintf("%s [%s]", msg, e.caller)
}

func (e *errorX) Wrap(err error) ErrorX {
	e.err = err

	return e
}

func (e *errorX) Skip(n int) ErrorX {
	e.skip = n
	e.setCaller(3)

	return e
}

func AsErrorX(err error) *errorX {
	if err == nil {
		return nil
	}
	if e, ok := err.(*errorX); ok {
		return e
	}
	return newf(err.Error()).(*errorX)
}

func New(message string) ErrorX {
	return newf(message)
}

func Newf(format string, args ...any) ErrorX {
	return newf(format, args...)
}

func NewHttp(status int, message string) ErrorX {
	return newHttpf(status, message)
}

func NewHttpf(status int, format string, args ...any) ErrorX {
	return newHttpf(status, format, args...)
}

func newHttpf(status int, format string, args ...any) ErrorX {
	e := &errorX{
		eType:   errorXHttpType,
		message: fmt.Sprintf(format, args...),

		errorXHttp: errorXHttp{
			status: status,
		},
	}
	e.setCaller(3)

	return e
}

func newf(format string, args ...any) ErrorX {
	e := &errorX{
		eType:   errorXDefaultType,
		message: fmt.Sprintf(format, args...),

		errorXHttp: errorXHttp{
			status: http.StatusInternalServerError,
		},
	}
	e.setCaller(3)

	return e
}

func (e *errorX) setCaller(skip int) {
	pc, _, line, _ := runtime.Caller(skip + e.skip)
	e.caller = fmt.Sprintf("%s:%d", runtime.FuncForPC(pc).Name(), line)
}
