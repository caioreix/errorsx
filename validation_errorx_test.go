package errorsx_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/caioreix/errorsx"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidationErrorX_Wrap(t *testing.T) {
	t.Parallel()
	type User struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	validate := validator.New()
	u1 := User{}
	validationErr1 := validate.Struct(u1)
	u2 := User{
		Name:  "x",
		Email: "x",
	}
	validationErr2 := validate.Struct(u2)

	var (
		msg = "foo"
	)

	rx := callerRX(fmt.Sprintf("%s: %s\n%s", msg, validationErr1.Error(), validationErr2.Error()))
	errX := errorsx.NewWithError(validationErr1, msg).Wrap(validationErr2)

	got := errX.Error()
	assert.Regexp(t, rx, got)
}

func TestValidationErrorX_Unwrap(t *testing.T) {
	t.Parallel()

	type User struct {
		Name  string `validate:"required"`
		Email string `validate:"required"`
	}

	validate := validator.New()
	u := User{}
	validationErr := validate.Struct(u)

	var (
		err = validationErr
		msg = "foo"
	)

	rx := callerRX(fmt.Sprintf("%s: %s", msg, err.Error()))
	errX := errorsx.NewWithError(err, msg)
	got := errX.Unwrap().Error()
	assert.Regexp(t, rx, got)
	assert.NotEqualValues(t, errX, got)
}

func TestValidationErrorX_NewWithError(t *testing.T) {
	type User struct {
		Name  string `validate:"required"`
		Email string `validate:"required"`
	}

	validate := validator.New()
	u := User{}
	validationErr := validate.Struct(u)

	t.Run("ErrorX", func(t *testing.T) {
		t.Parallel()
		var (
			err = validationErr
			msg = "foo"
		)

		rx := callerRX(fmt.Sprintf("%s: %s", msg, err.Error()))
		errX := errorsx.NewWithError(err, msg)

		got := errX.Error()
		assert.Regexp(t, rx, got)
	})

	t.Run("HTTPErrorX", func(t *testing.T) {
		t.Parallel()
		var (
			err    = validationErr
			msg    = "foo"
			status = http.StatusBadRequest
		)

		rx := callerRX(fmt.Sprintf("%s: %s: status %d", msg, err.Error(), status))
		errX := errorsx.NewHTTPWithError(err, status, msg)

		got := errX.Error()
		assert.Regexp(t, rx, got)
	})
}

func TestValidationErrorX_NewWithWrap(t *testing.T) {
	type User struct {
		Name  string `validate:"required"`
		Email string `validate:"required"`
	}

	validate := validator.New()
	u := User{}
	validationErr := validate.Struct(u)

	t.Run("ErrorX", func(t *testing.T) {
		t.Parallel()
		var (
			err = validationErr
			msg = "foo"
		)

		rx := callerRX(fmt.Sprintf("%s: %s", msg, err.Error()))
		errX := errorsx.New(msg).Wrap(err)

		got := errX.Error()
		assert.Regexp(t, rx, got)
	})

	t.Run("HTTPErrorX", func(t *testing.T) {
		t.Parallel()
		var (
			err    = validationErr
			msg    = "foo"
			status = http.StatusBadRequest
		)

		rx := callerRX(fmt.Sprintf("%s: %s: status %d", msg, err.Error(), status))
		errX := errorsx.NewHTTP(status, msg).Wrap(err)

		got := errX.Error()
		assert.Regexp(t, rx, got)
	})
}

func TestValidationErrorX_Fields(t *testing.T) {
	type User struct {
		Name  string `validate:"required"`
		Email string `validate:"required,email"`
	}

	validate := validator.New()
	u := User{
		Email: "x",
	}
	validationErr := validate.Struct(u)

	t.Run("ErrorX", func(t *testing.T) {
		t.Parallel()
		var (
			err = validationErr
			msg = "foo"

			want = map[string]any{
				"message": msg,
				"error":   err.Error(),
				"validation_errors": map[string]any{
					"Email": "email",
					"Name":  "required",
				},
			}
		)

		errX := errorsx.NewWithError(err, msg)
		want["caller"] = errX.Caller()
		want["stack"] = errX.Stack()

		got := errX.Fields()
		assert.Equal(t, want, got)
	})

	t.Run("HTTPErrorX", func(t *testing.T) {
		t.Parallel()
		var (
			err    = validationErr
			msg    = "foo"
			status = http.StatusBadRequest

			want = map[string]any{
				"message": msg,
				"error":   err.Error(),
				"status":  status,
				"validation_errors": map[string]any{
					"Email": "email",
					"Name":  "required",
				},
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
			err = validationErr
			msg = "foo"

			want = map[string]any{
				"message": msg,
				"validation_errors": map[string]any{
					"Email": "email",
					"Name":  "required",
				},
			}
		)

		errX := errorsx.NewWithError(err, msg)
		got := errX.Fields("message", "validation_errors")
		assert.Equal(t, want, got)
	})

	t.Run("nil ValidationErrors", func(t *testing.T) {
		t.Parallel()
		var (
			err validator.ValidationErrors = nil
			msg                            = "foo"

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

	t.Run("empty ValidationErrors", func(t *testing.T) {
		t.Parallel()
		var (
			err = validator.ValidationErrors{}
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
}
