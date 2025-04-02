package errorsx

import (
	"strconv"
)

type httpErrorX struct {
	ErrorX

	status int
}

func (e *httpErrorX) Error() string {
	return stringify(e)
}

func (e *httpErrorX) Unwrap() error {
	return e.unwrap()
}

func (e httpErrorX) Wrap(err error) ErrorX {
	e.ErrorX = e.ErrorX.Wrap(err)
	return &e
}

func (e *httpErrorX) string() string {
	return "status " + strconv.FormatInt(int64(e.status), 10)
}

func (e *httpErrorX) Fields(fields ...string) map[string]any {
	return mapify(e, fields)
}

func (e httpErrorX) unwrap() ErrorX {
	return e.ErrorX
}

func (e *httpErrorX) fields() map[string]any {
	return map[string]any{"status": e.status}
}

func NewHTTP(status int, message string) ErrorX {
	return &httpErrorX{
		ErrorX: newf(nil, "%s", message),
		status: status,
	}
}

func NewHTTPf(status int, format string, args ...any) ErrorX {
	return &httpErrorX{
		ErrorX: newf(nil, format, args...),
		status: status,
	}
}

func NewHTTPWithError(err error, status int, message string) ErrorX {
	return &httpErrorX{
		ErrorX: newf(err, "%s", message),
		status: status,
	}
}

func NewHTTPWithErrorf(err error, status int, format string, args ...any) ErrorX {
	return &httpErrorX{
		ErrorX: newf(err, format, args...),
		status: status,
	}
}
