package utils

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// Returns a status code between 300 and 599
// {"status": "text for the status code", "message": "the error message"}
func Error(w http.ResponseWriter, statusCode int, message string) {
	defaultStatusCode := http.StatusInternalServerError
	// check if status code is valid
	if statusCode > 299 && statusCode < 600 {
		defaultStatusCode = statusCode
	}

	body := errorResponse{
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
	w.Write(bytes)
}

type successResponse struct {
	Data any `json:"data"`
}

// Returns a success status.
// e.g {"data": [...]}
func JSON(w http.ResponseWriter, code int, body any) {
	// check body
	if body == nil {
		w.WriteHeader(code)
		return
	}
	var bytes []byte
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
	w.Write(bytes)
}
