package errorsx_test

import (
	"errors"
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

func TestErrorX_Fields(t *testing.T) {
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

		got := errX.Fields()
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
		got := errX.Fields("message", "status")
		assert.Equal(t, want, got)
	})
}

func TestErrorX_Wrap(t *testing.T) {
	t.Parallel()

	t.Run("wrap first error", func(t *testing.T) {
		t.Parallel()
		msg := "original error"
		wrapErr := fmt.Errorf("wrapped error")

		errX := errorsx.New(msg)
		wrappedErrX := errX.Wrap(wrapErr)

		expected := fmt.Sprintf("%s: %s", msg, wrapErr.Error())
		rx := callerRX(expected)

		got := wrappedErrX.Error()
		assert.Regexp(t, rx, got)

		resp := wrappedErrX.Fields()
		assert.Contains(t, resp["error"].(string), wrapErr.Error())
	})

	t.Run("wrap multiple errors", func(t *testing.T) {
		t.Parallel()
		msg := "original error"
		firstWrapErr := fmt.Errorf("first wrapped error")
		secondWrapErr := fmt.Errorf("second wrapped error")

		errX := errorsx.New(msg)
		wrappedOnce := errX.Wrap(firstWrapErr)
		wrappedTwice := wrappedOnce.Wrap(secondWrapErr)

		got := wrappedTwice.Error()
		assert.Contains(t, got, msg)
		assert.Contains(t, got, firstWrapErr.Error())
		assert.Contains(t, got, secondWrapErr.Error())

		resp := wrappedTwice.Fields()
		assert.Contains(t, resp["error"].(string), firstWrapErr.Error())
		assert.Contains(t, resp["error"].(string), secondWrapErr.Error())
	})

	t.Run("wrap nil error", func(t *testing.T) {
		t.Parallel()
		msg := "original error"

		errX := errorsx.New(msg)
		wrappedErrX := errX.Wrap(nil)

		got := wrappedErrX.Error()
		assert.Contains(t, got, msg)

		resp := wrappedErrX.Fields()
		assert.Equal(t, msg, resp["message"])
	})
}

func TestErrorX_Unwrap(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name     string
		setupErr func() (errorsx.ErrorX, error)
		wantErr  error
	}{
		{
			name: "unwrap with wrapped error",
			setupErr: func() (errorsx.ErrorX, error) {
				wrappedErr := fmt.Errorf("wrapped error")
				return errorsx.NewWithError(wrappedErr, "parent error"), wrappedErr
			},
			wantErr: fmt.Errorf("wrapped error"),
		},
		{
			name: "unwrap with nil error",
			setupErr: func() (errorsx.ErrorX, error) {
				return errorsx.New("error without wrapped error"), nil
			},
			wantErr: nil,
		},
		{
			name: "unwrap with multiple errors",
			setupErr: func() (errorsx.ErrorX, error) {
				originalErr := fmt.Errorf("original error")
				errX := errorsx.NewWithError(originalErr, "parent error")
				additionalErr := fmt.Errorf("additional error")
				errX = errX.Wrap(additionalErr)
				return errX, errors.Join(originalErr, additionalErr)
			},
			wantErr: errors.Join(fmt.Errorf("original error"), fmt.Errorf("additional error")),
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			errX, wantErr := tc.setupErr()

			gotErr := errX.Unwrap()

			if wantErr == nil {
				assert.Nil(t, gotErr)
			} else {
				assert.Equal(t, wantErr.Error(), gotErr.Error())
			}
		})
	}
}

func callerRX(msg string, skip ...int) string {
	skipT := 0
	for _, s := range skip {
		skipT += s
	}

	pc, file, _, _ := runtime.Caller(skipT + 1)
	return fmt.Sprintf(`^%s \[%s %s:\d+\]$`, msg, runtime.FuncForPC(pc).Name(), file)
}
