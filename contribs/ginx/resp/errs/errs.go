package errs

import "net/http"

// Error represents http response error, which is along with http status code,
// it used to decide how to show error message in response.
type Error struct {
	Err    error
	Status int
}

func (e Error) Error() string {
	return e.Err.Error()
}

func Wrap(err error, status int) Error {
	return Error{Err: err, Status: status}
}

func BadRequest(err error) Error {
	return Wrap(err, http.StatusBadRequest)
}

func Internal(err error) Error {
	return Wrap(err, http.StatusInternalServerError)
}

func UnAuthorized(err error) Error {
	return Wrap(err, http.StatusUnauthorized)
}

func Forbidden(err error) Error {
	return Wrap(err, http.StatusForbidden)
}
