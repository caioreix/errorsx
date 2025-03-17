package errorsx_test

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/caioreix/errorsx"
)

func TestErrorX_Error(t *testing.T) {
	tt := []struct {
		name string
		err  errorsx.ErrorX
		want string
	}{
		{
			name: "without wrap",
			err:  errorsx.New("foo"),
			want: "foo",
		},
		{
			name: "with wrap",
			err:  errorsx.New("foo").Wrap(errors.New("bar")),
			want: "foo: bar",
		},
	}
	rx := callerRX("%s")

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.err.Error()
			assert.Regexp(t, fmt.Sprintf(rx, tc.want), got)
		})
	}
}

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

func TestErrorX_NewHttp(t *testing.T) {
	t.Parallel()
	var (
		msg    = "foo"
		status = http.StatusUnprocessableEntity
	)

	rx := callerRX(fmt.Sprintf("%s: status: %d", msg, status))
	err := errorsx.NewHttp(status, msg)
	got := err.Error()
	assert.Regexp(t, rx, got)
}

func TestErrorX_NewHttpf(t *testing.T) {
	t.Parallel()
	var (
		format = "foo %s"
		args   = []any{"bar"}
		status = http.StatusUnprocessableEntity
	)

	rx := callerRX(fmt.Sprintf("%s: status: %d", fmt.Sprintf(format, args...), status))
	err := errorsx.NewHttpf(status, format, args...)
	got := err.Error()
	assert.Regexp(t, rx, got)
}

func callerRX(msg string, skip ...int) string {
	skipT := 0
	for _, s := range skip {
		skipT += s
	}

	pc, _, _, _ := runtime.Caller(skipT + 1)
	return fmt.Sprintf(`^%s \[%s:\d+\]$`, msg, runtime.FuncForPC(pc).Name())
}
