package utils

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test__JSON__WhenSlice(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, 204, []int{1, 2, 3})

	bContent, _ := io.ReadAll(w.Result().Body)
	resultBody := strings.TrimSpace(string(bContent))

	require.Equal(t, w.Result().StatusCode, 204)
	require.Equal(t, `{"data":[1,2,3]}`, resultBody)
	require.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
}

func Test__JSON__WhenSimpleData(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, 204, "Some simple data")

	bContent, _ := io.ReadAll(w.Result().Body)
	response := strings.TrimSpace(string(bContent))

	require.Equal(t, w.Result().StatusCode, 204)
	require.Equal(t, `{"data":"Some simple data"}`, response)
	require.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
}

func Test__JSON__WhenNil(t *testing.T) {
	w := httptest.NewRecorder()
	JSON(w, 204, nil)

	bContent, _ := io.ReadAll(w.Result().Body)
	response := strings.TrimSpace(string(bContent))

	require.Equal(t, w.Result().StatusCode, 204)
	require.Equal(t, "", response)
}

func Test__Error(t *testing.T) {
	w := httptest.NewRecorder()
	Error(w, 500, "Some error occurs")

	var resultBody ErrorResponse
	json.NewDecoder(w.Result().Body).Decode(&resultBody)

	require.Equal(t, w.Result().StatusCode, 500)
	require.Equal(t, resultBody, ErrorResponse{"Internal Server Error", "Some error occurs"})
	require.Equal(t, w.Result().Header.Get("Content-Type"), "application/json")
}
