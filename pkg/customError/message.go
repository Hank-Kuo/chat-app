package customError

import "errors"

var (
	ErrBadRequest             = errors.New("Bad request")
	ErrWrongCredentials       = errors.New("Wrong Credentials")
	ErrNotFound               = errors.New("Not Found")
	ErrUnauthorized           = errors.New("Unauthorized")
	ErrForbidden              = errors.New("Forbidden")
	ErrExpiredCSRFError       = errors.New("Expired CSRF token")
	ErrWrongCSRFToken         = errors.New("Wrong CSRF token")
	ErrCSRFNotPresented       = errors.New("CSRF not presented")
	ErrNotRequiredFields      = errors.New("No such required fields")
	ErrBadQueryParams         = errors.New("Invalid query params")
	ErrInternalServerError    = errors.New("Internal Server Error")
	ErrRequestTimeoutError    = errors.New("Request Timeout")
	ErrNotGetJWTToken         = errors.New("Not Get JWT claims")
	ErrInvalidJWTClaims       = errors.New("Invalid JWT claims")
	ErrExpiredJWTError        = errors.New("JWT claims is expired")
	ErrNotAllowedImageHeader  = errors.New("Not allowed image header")
	ErrContextCancel          = errors.New("Context cacnel")
	ErrDeadlineExceeded       = errors.New("Context exceed")
	ErrPasswordCodeNotMatched = errors.New("Password is wrong")
)
