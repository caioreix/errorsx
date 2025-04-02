package errorsx

import (
	"errors"
	"fmt"
	"runtime"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ErrorX interface {
	error

	Caller() string
	Stack() Stack
	Fields(fields ...string) map[string]any
	Wrap(err error) ErrorX
	Unwrap() error

	// helpers
	string() string
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

func (e *errorX) string() string {
	msg := e.message
	if e.err != nil {
		msg = msg + ": " + e.err.Error()
	}

	return msg
}

func (e errorX) Wrap(err error) ErrorX {
	e.err = errors.Join(e.err, err)
	switch et := err.(type) {
	case validator.ValidationErrors:
		return &validationErrorX{
			ErrorX:      &e,
			fieldErrors: et,
		}
	}

	return &e
}

func (e *errorX) Unwrap() error {
	return e.err
}

func (e *errorX) Fields(fields ...string) map[string]any {
	return mapify(e, fields)
}

func (e *errorX) Caller() string {
	return e.caller
}

func (e *errorX) Stack() Stack {
	return e.stack
}

func (e errorX) unwrap() ErrorX {
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
	newErrorX := &errorX{
		err:     err,
		message: fmt.Sprintf(format, args...),
		caller:  getCaller(2),
		stack:   getStack(4),
	}

	switch e := err.(type) {
	case validator.ValidationErrors:
		return &validationErrorX{
			ErrorX:      newErrorX,
			fieldErrors: e,
		}
	default:
		return newErrorX
	}
}

func stringify(e ErrorX) string {
	ex := ErrorX(e)
	msgs := make([]string, 1)
	for {
		if eu := ex.unwrap(); eu != nil {
			msg := ex.string()
			if msg != "" {
				msgs = append(msgs, msg)
			}
			ex = eu
			continue
		}

		msg := ex.string()

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

		if len(fields) == 0 || slices.Contains(fields, "stack") {
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
