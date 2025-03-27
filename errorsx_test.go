package errorsx_test

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/caioreix/errorsx"
	"github.com/stretchr/testify/assert"
)

func TestErrorX_New(t *testing.T) {
	t.Parallel()
	var msg = "foo"

	rx := callerRX(msg)
	err := errorsx.New(msg)
	got := err.Error()
	assert.Regexp(t, rx, got)
}

func TestErrorX_Newf(t *testing.T) {
	t.Parallel()
	var (
		format = "foo %s"
		args   = []any{"bar"}
	)

	rx := callerRX(fmt.Sprintf(format, args...))
	err := errorsx.Newf(format, args...)
	got := err.Error()
	assert.Regexp(t, rx, got)
}

func TestErrorX_NewWithError(t *testing.T) {
	t.Parallel()
	var (
		err = fmt.Errorf("fake error")
		msg = "foo"
	)

	rx := callerRX(fmt.Sprintf("%s: %s", msg, err.Error()))
	errX := errorsx.NewWithError(err, msg)
	got := errX.Error()
	assert.Regexp(t, rx, got)
}

func TestErrorX_NewWithErrorf(t *testing.T) {
	t.Parallel()
	var (
		err    = fmt.Errorf("fake error")
		format = "foo %s"
		args   = []any{"bar"}
	)

	rx := callerRX(fmt.Sprintf("%s: %s", fmt.Sprintf(format, args...), err.Error()))
	errX := errorsx.NewWithErrorf(err, format, args...)
	got := errX.Error()
	assert.Regexp(t, rx, got)
}

func TestErrorX_Response(t *testing.T) {
	t.Run("without filter", func(t *testing.T) {
		t.Parallel()
		var (
			err = fmt.Errorf("fake error")
			msg = "foo"

			want = map[string]any{
				"message": msg,
				"error":   err.Error(),
			}
		)

		errX := errorsx.NewWithError(err, msg)
		want["caller"] = errX.Caller()
		want["stack"] = errX.Stack()

		got := errX.Response()
		assert.Equal(t, want, got)
	})

	t.Run("with filter", func(t *testing.T) {
		t.Parallel()
		var (
			err = fmt.Errorf("fake error")
			msg = "foo"

			want = map[string]any{
				"message": msg,
			}
		)

		errX := errorsx.NewWithError(err, msg)
		got := errX.Response("message", "status")
		assert.Equal(t, want, got)
	})
}

func callerRX(msg string, skip ...int) string {
	skipT := 0
	for _, s := range skip {
		skipT += s
	}

	pc, file, _, _ := runtime.Caller(skipT + 1)
	return fmt.Sprintf(`^%s \[%s %s:\d+\]$`, msg, runtime.FuncForPC(pc).Name(), file)
}
