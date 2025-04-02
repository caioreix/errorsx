package errorsx

import "github.com/go-playground/validator/v10"

type validationErrorX struct {
	ErrorX

	fieldErrors validator.ValidationErrors
}

func (e *validationErrorX) Error() string {
	return stringify(e)
}

func (e *validationErrorX) Unwrap() error {
	return e.unwrap()
}

func (e validationErrorX) Wrap(err error) ErrorX {
	e.ErrorX = e.ErrorX.Wrap(err)
	return &e
}

func (e *validationErrorX) string() string {
	return ""
}

func (e *validationErrorX) Fields(fields ...string) map[string]any {
	return mapify(e, fields)
}

func (e validationErrorX) unwrap() ErrorX {
	return e.ErrorX
}

func (e *validationErrorX) fields() map[string]any {
	m := map[string]any{}

	if e.fieldErrors == nil {
		return m
	}

	errs := make(map[string]any)
	for _, err := range e.fieldErrors {
		errs[err.Field()] = err.Tag()
	}

	if len(errs) != 0 {
		m["validation_errors"] = errs
	}

	return m
}
