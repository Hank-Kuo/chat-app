package httpError

import (
	"fmt"
	"net/http"

	"github.com/Hank-Kuo/chat-app/pkg/customError"
)

type Err struct {
	Status  int
	Message string
	Detail  interface{}
}

func (e Err) Error() string {
	return fmt.Sprintf("%v", e.Message)
}

func (e Err) GetMessage() string {
	if e.Message == "" {
		return "Internal server error"
	}
	return e.Message
}

func (e Err) GetStatus() int {
	if e.Status == 0 {
		return http.StatusInternalServerError
	}
	return e.Status
}

func NewError(httpCode int, message string, err error) *Err {
	return &Err{
		httpCode, message, err,
	}
}

func NewNotFoundError(err error) *Err {
	return &Err{
		http.StatusNotFound, customError.ErrNotFound.Error(), err,
	}
}

func NewInternalServerError(err error) *Err {
	return &Err{
		http.StatusInternalServerError, customError.ErrInternalServerError.Error(), err,
	}
}

func NewBadRequestError(err error) *Err {
	return &Err{
		http.StatusBadRequest, customError.ErrBadRequest.Error(), err,
	}
}

func NewUnauthorizedError(err error) *Err {
	return &Err{
		http.StatusUnauthorized, customError.ErrUnauthorized.Error(), err,
	}
}
func NewRequestTimeoutError(err error) *Err {
	return &Err{
		http.StatusRequestTimeout, customError.ErrRequestTimeoutError.Error(), err,
	}
}

func NewBadQueryParamsError(err error) *Err {
	return &Err{
		http.StatusBadRequest, customError.ErrBadQueryParams.Error(), err,
	}
}
