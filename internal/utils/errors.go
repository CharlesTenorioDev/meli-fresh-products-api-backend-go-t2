package utils

import (
	"errors"
	"net/http"

	"github.com/bootcamp-go/web/response"
)

var (
	ErrInvalidFormat          = errors.New("invalid format")               // 400
	ErrInvalidArguments       = errors.New("invalid arguments")            // 422
	ErrConflict               = errors.New("entity already exists")        // 409
	ErrNotFound               = errors.New("entity not found")             // 404
	ErrInvalidProperties      = errors.New("invalid properties format")    // For parsing the properties, panic
	ErrEmptyArguments         = errors.New("arguments must not be empty")  // 422
	ErrWarehouseDoesNotExists = errors.New("warehouse's id doesn't exist") // 422
	ErrBuyerDoesNotExists     = errors.New("buyer's id doesn't exist")     // 409
	ErrProductDoesNotExists   = errors.New("product's id doesn't exist")   // 409
)

func ENotFound(target string) error {
	return errors.Join(ErrNotFound, errors.New(target+" doesn't exist"))
}

func EZeroValue(target string) error {
	return errors.Join(ErrNotFound, errors.New(target+" cannot be empty/null"))
}

func EConflict(attribute, target string) error {
	return errors.Join(ErrInvalidArguments, errors.New(target+" with attribute '"+attribute+"' alredy exists"))
}

func EDependencyNotFound(attribute, target string) error {
	return errors.Join(ErrInvalidArguments, errors.New(target+" with '"+attribute+"' doesn't exist"))
}
func EBR(message string) error {
	return errors.Join(ErrInvalidArguments, errors.New(message))
}

// HandleError centralizes error handling and response formatting
func HandleError(w http.ResponseWriter, err error) {
	var status int

	var message string

	switch err {
	case ErrInvalidFormat:
		status = http.StatusBadRequest
		message = ErrInvalidFormat.Error()
	case ErrWarehouseDoesNotExists:
		status = http.StatusUnprocessableEntity
		message = ErrWarehouseDoesNotExists.Error()
	case ErrInvalidArguments:
		status = http.StatusUnprocessableEntity
		message = ErrInvalidArguments.Error()
	case ErrEmptyArguments:
		status = http.StatusUnprocessableEntity
		message = ErrEmptyArguments.Error()
	case ErrConflict:
		status = http.StatusConflict
		message = ErrConflict.Error()
	case ErrNotFound:
		status = http.StatusNotFound
		message = ErrNotFound.Error()
	case ErrInvalidProperties:
		status = http.StatusBadRequest
		message = ErrInvalidProperties.Error()
	case ErrBuyerDoesNotExists:
		status = http.StatusConflict
		message = ErrBuyerDoesNotExists.Error()
	case ErrProductDoesNotExists:
		status = http.StatusConflict
		message = ErrProductDoesNotExists.Error()
	default:
		status = http.StatusInternalServerError
		message = "internal server error"
	}

	response.Error(w, status, message)
}
