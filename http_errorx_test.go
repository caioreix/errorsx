package errorsx_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/caioreix/errorsx"
	"github.com/stretchr/testify/assert"
)

func TestErrorX_NewHttp(t *testing.T) {
	t.Parallel()
	var (
		msg    = "foo"
		status = http.StatusUnprocessableEntity
	)

	rx := callerRX(fmt.Sprintf("%s: status %d", msg, status))
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

	rx := callerRX(fmt.Sprintf("%s: status %d", fmt.Sprintf(format, args...), status))
	err := errorsx.NewHttpf(status, format, args...)
	got := err.Error()
	assert.Regexp(t, rx, got)
}

func TestErrorX_NewHttpWithError(t *testing.T) {
	t.Parallel()
	var (
		err    = fmt.Errorf("fake error")
		msg    = "foo"
		status = http.StatusUnprocessableEntity
	)

	rx := callerRX(fmt.Sprintf("%s: %s: status %d", msg, err.Error(), status))
	errX := errorsx.NewHttpWithError(err, status, msg)
	got := errX.Error()
	assert.Regexp(t, rx, got)
}

func TestErrorX_NewHttpWithErrorf(t *testing.T) {
	t.Parallel()
	var (
		err    = fmt.Errorf("fake error")
		format = "foo %s"
		args   = []any{"bar"}
		status = http.StatusUnprocessableEntity
	)

	rx := callerRX(fmt.Sprintf("%s: %s: status %d", fmt.Sprintf(format, args...), err.Error(), status))
	errX := errorsx.NewHttpWithErrorf(err, status, format, args...)
	got := errX.Error()
	assert.Regexp(t, rx, got)
}
