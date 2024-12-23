package utils

import (
	"errors"
)

var (
	ErrInvalidFormat     = errors.New("invalid format")            // 400
	ErrInvalidArguments  = errors.New("invalid arguments")         // 422
	ErrConflict          = errors.New("entity already exists")     // 409
	ErrNotFound          = errors.New("entity not found")          // 404
	ErrInvalidProperties = errors.New("invalid properties format") // For parsing the properties, panic
)
