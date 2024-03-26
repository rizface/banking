package db

import "errors"

var (
	ErrNoRow                = errors.New("data not found")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrProductNameDuplicate = errors.New("product name already exists")
)
