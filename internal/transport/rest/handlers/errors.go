package handlers

import "errors"

var (
	ErrInternalServerError = errors.New("internal Server Error")
	ErrBadRequest          = errors.New("request is invalid")
	ErrUnauthorized        = errors.New("request is unauthorized")
)
