package rest

import (
	"net/http"
)

//var (
//	//errUserNotFound       = errors.New("wrong login or password")
//	errUserExist          = errors.New("user exist already")
//	errEmptyCookie        = errors.New("cookie is empty")
//	errEmptyCookieAccess  = errors.New("access cookie is empty")
//	errEmptyCookieRefresh = errors.New("access cookie is empty")
//)

// swagger:model Error
type Error struct {
	Error string `json:"error"`
}

type HttpError struct {
	StatusCode int
	Message    error
}

func (he HttpError) Error() string {
	return he.Message.Error()
}

func (he HttpError) Status() int {
	return he.StatusCode
}

func ErrorBadRequest(err error) HttpError {
	return HttpError{
		StatusCode: http.StatusBadRequest,
		Message:    err,
	}
}

func ErrorForbidden(err error) HttpError {
	return HttpError{
		StatusCode: http.StatusForbidden,
		Message:    err,
	}
}

func ErrorNotFound(err error) HttpError {
	return HttpError{
		StatusCode: http.StatusNotFound,
		Message:    err,
	}
}

func ErrorInternal(err error) HttpError {
	return HttpError{
		StatusCode: http.StatusInternalServerError,
		Message:    err,
	}
}

func ErrConflict(err error) HttpError {
	return HttpError{
		StatusCode: http.StatusConflict,
		Message:    err,
	}
}
