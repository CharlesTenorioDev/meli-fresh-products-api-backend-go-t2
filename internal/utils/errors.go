package utils

import (
	"errors"
	"net/http"
	"strings"

	"github.com/bootcamp-go/web/response"
)

var (
	ErrInvalidFormat     = errors.New("invalid format")              // 400
	ErrInvalidArguments  = errors.New("invalid arguments")           // 422
	ErrConflict          = errors.New("entity already exists")       // 409
	ErrNotFound          = errors.New("entity not found")            // 404
	ErrInvalidProperties = errors.New("invalid properties format")   // For parsing the properties, panic
	ErrEmptyArguments    = errors.New("arguments must not be empty") // 422
)

// ENotFound When 404 status, only when entity has some relation with the url, e.g. GET /product/1
// and no product exist for 1
func ENotFound(target string) error {
	return errors.Join(ErrNotFound, errors.New(target+" doesn't exist"))
}

// EZeroValue When 422 status, when attribute is zero value (empty, nil, 0, {})
func EZeroValue(target string) error {
	return errors.Join(ErrInvalidArguments, errors.New(target+" cannot be empty/null"))
}

// EConflict When 409, when trying to manage a resource, and some attribute already exist
func EConflict(attribute, target string) error {
	return errors.Join(ErrConflict, errors.New(target+" with attribute '"+attribute+"' already exists"))
}

// EDependencyNotFound When 422, when managing a resource, and an attribute that refers other entity
// does not exist for that value
func EDependencyNotFound(target, attribute string) error {
	return errors.Join(ErrInvalidArguments, errors.New(target+" with '"+attribute+"' doesn't exist"))
}

// EBR When 422, for business rule validation
func EBR(message string) error {
	return errors.Join(ErrInvalidArguments, errors.New(message))
}

// EBadRequest When 400, when payload or query params or path value cannot be processed
// due to their format
func EBadRequest(attribute string) error {
	return errors.Join(ErrInvalidFormat, errors.New(attribute+" with invalid format"))
}

// HandleError centralizes error handling and response formatting
func HandleError(w http.ResponseWriter, err error) {
	var status int

	var message string

	if err == nil {
		status = http.StatusInternalServerError
		message = "internal server error"
	}

	if errors.Is(err, ErrInvalidFormat) {
		status = http.StatusBadRequest
		message = err.Error()
	} else if errors.Is(err, ErrInvalidArguments) {
		status = http.StatusUnprocessableEntity
		message = err.Error()
	} else if errors.Is(err, ErrEmptyArguments) {
		status = http.StatusUnprocessableEntity
		message = err.Error()
	} else if errors.Is(err, ErrConflict) {
		status = http.StatusConflict
		message = err.Error()
	} else if errors.Is(err, ErrNotFound) {
		status = http.StatusNotFound
		message = err.Error()
	} else {
		status = http.StatusInternalServerError
		message = "internal server error"
	}

	message = strings.Replace(message, "\n", ": ", 1)
	response.Error(w, status, message)
}
