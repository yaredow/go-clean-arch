package domain

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrNotFound            = errors.New("your request item not found")
	ErrConflict            = errors.New("your item already exists")
	ErrBadParamInput       = errors.New("given param is not valid")
)
