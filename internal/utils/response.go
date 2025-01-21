package utils

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Error returns a status code between 300 and 599
// {"status": "text for the status code", "message": "the error message"}
func Error(w http.ResponseWriter, statusCode int, message string) {
	defaultStatusCode := http.StatusInternalServerError
	// check if status code is valid
	if statusCode > 299 && statusCode < 600 {
		defaultStatusCode = statusCode
	}

	body := ErrorResponse{
		Status:  http.StatusText(defaultStatusCode),
		Message: message,
	}
	bytes, err := json.Marshal(body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// write response
	// - set header: before code due to it sets by default "text/plain"
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(defaultStatusCode)

	if _, err = w.Write(bytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type successResponse struct {
	Data any `json:"data"`
}

func JSON(w http.ResponseWriter, code int, body any) {
	// check body
	if body == nil {
		w.WriteHeader(code)
		return
	}

	response := successResponse{body}
	bytes, err := json.Marshal(response)

	if err != nil {
		// default error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// set header (before code due to it sets by default "text/plain")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if _, err = w.Write(bytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
