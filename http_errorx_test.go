package errorsx_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/caioreix/errorsx"
	"github.com/stretchr/testify/assert"
)

func TestHTTPErrorX_Wrap(t *testing.T) {
	t.Parallel()
	var (
		err    = fmt.Errorf("fake error")
		msg    = "foo"
		status = http.StatusUnprocessableEntity
	)

	rx := callerRX(fmt.Sprintf("%s: %s: status %d", msg, err.Error(), status))
	errX := errorsx.NewHTTP(status, msg).Wrap(err)
	got := errX.Error()
	assert.Regexp(t, rx, got)
}

func TestHTTPErrorX_Unwrap(t *testing.T) {
	t.Parallel()
	var (
		err    = fmt.Errorf("fake error")
		msg    = "foo"
		status = http.StatusUnprocessableEntity
	)

	rx := callerRX(fmt.Sprintf("%s: %s", msg, err.Error()))
	errX := errorsx.NewHTTPWithError(err, status, msg)
	got := errX.Unwrap().Error()
	assert.Regexp(t, rx, got)
}

func TestHTTPErrorX_Fields(t *testing.T) {
	t.Run("without filter", func(t *testing.T) {
		t.Parallel()
		var (
			err    = fmt.Errorf("fake error")
			msg    = "foo"
			status = http.StatusUnprocessableEntity

			want = map[string]any{
				"message": msg,
				"status":  status,
				"error":   err.Error(),
			}
		)

		errX := errorsx.NewHTTPWithError(err, status, msg)
		want["caller"] = errX.Caller()
		want["stack"] = errX.Stack()

		got := errX.Fields()
		assert.Equal(t, want, got)
	})

	t.Run("with filter", func(t *testing.T) {
		t.Parallel()
		var (
			err    = fmt.Errorf("fake error")
			msg    = "foo"
			status = http.StatusUnprocessableEntity

			want = map[string]any{
				"message": msg,
				"status":  status,
			}
		)

		errX := errorsx.NewHTTPWithError(err, status, msg)
		got := errX.Fields("message", "status")
		assert.Equal(t, want, got)
	})
}

func TestHTTPErrorX_NewHTTP(t *testing.T) {
	t.Parallel()
	var (
		msg    = "foo"
		status = http.StatusUnprocessableEntity
	)

	rx := callerRX(fmt.Sprintf("%s: status %d", msg, status))
	err := errorsx.NewHTTP(status, msg)
	got := err.Error()
	assert.Regexp(t, rx, got)
}

func TestHTTPErrorX_NewHTTPf(t *testing.T) {
	t.Parallel()
	var (
		format = "foo %s"
		args   = []any{"bar"}
		status = http.StatusUnprocessableEntity
	)

	rx := callerRX(fmt.Sprintf("%s: status %d", fmt.Sprintf(format, args...), status))
	err := errorsx.NewHTTPf(status, format, args...)
	got := err.Error()
	assert.Regexp(t, rx, got)
}

func TestHTTPErrorX_NewHTTPWithError(t *testing.T) {
	t.Parallel()
	var (
		err    = fmt.Errorf("fake error")
		msg    = "foo"
		status = http.StatusUnprocessableEntity
	)

	rx := callerRX(fmt.Sprintf("%s: %s: status %d", msg, err.Error(), status))
	errX := errorsx.NewHTTPWithError(err, status, msg)
	got := errX.Error()
	assert.Regexp(t, rx, got)
}

func TestHTTPErrorX_NewHTTPWithErrorf(t *testing.T) {
	t.Parallel()
	var (
		err    = fmt.Errorf("fake error")
		format = "foo %s"
		args   = []any{"bar"}
		status = http.StatusUnprocessableEntity
	)

	rx := callerRX(fmt.Sprintf("%s: %s: status %d", fmt.Sprintf(format, args...), err.Error(), status))
	errX := errorsx.NewHTTPWithErrorf(err, status, format, args...)
	got := errX.Error()
	assert.Regexp(t, rx, got)
}
