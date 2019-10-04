package handlers

type HTTPError struct {
	err  error
	Code int
}

func (e HTTPError) Error() string {
	return e.err.Error()
}

func NewHTTPError(err error, code int) HTTPError {
	return HTTPError{
		err:  err,
		Code: code,
	}
}
