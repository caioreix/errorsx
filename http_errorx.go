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

func (e *httpErrorX) String() string {
	return "status " + strconv.FormatInt(int64(e.status), 10)
}

func (e *httpErrorX) unwrap() ErrorX {
	return e.ErrorX
}

func NewHttp(status int, message string) ErrorX {
	return &httpErrorX{
		ErrorX: newf(nil, "%s", message),
		status: status,
	}
}

func NewHttpf(status int, format string, args ...any) ErrorX {
	return &httpErrorX{
		ErrorX: newf(nil, format, args...),
		status: status,
	}
}

func NewHttpWithError(err error, status int, message string) ErrorX {
	return &httpErrorX{
		ErrorX: newf(err, "%s", message),
		status: status,
	}
}

func NewHttpWithErrorf(err error, status int, format string, args ...any) ErrorX {
	return &httpErrorX{
		ErrorX: newf(err, format, args...),
		status: status,
	}
}
