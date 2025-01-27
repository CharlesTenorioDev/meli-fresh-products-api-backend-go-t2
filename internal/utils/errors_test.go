package utils

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestENotFound(t *testing.T) {
	err := ENotFound("some_entity")
	require.ErrorIs(t, err, ErrNotFound)
	require.Equal(t, err.Error(), "entity not found\nsome_entity doesn't exist")
}
func TestEZeroValue(t *testing.T) {
	err := EZeroValue("some_attribute")
	require.ErrorIs(t, err, ErrInvalidArguments)
	require.Equal(t, err.Error(), "invalid arguments\nsome_attribute cannot be empty/null")
}
func TestEConflict(t *testing.T) {
	err := EConflict("foo", "bar")
	require.ErrorIs(t, err, ErrConflict)
	require.Equal(t, err.Error(), "entity already exists\nfoo with attribute 'bar' already exists")
}
func TestEDependencyNotFound(t *testing.T) {
	err := EDependencyNotFound("foo", "bar")
	require.ErrorIs(t, err, ErrInvalidArguments)
	require.Equal(t, err.Error(), "invalid arguments\nfoo with 'bar' doesn't exist")
}
func TestEBR(t *testing.T) {
	err := EBR("foo bar is invalid")
	require.ErrorIs(t, err, ErrInvalidArguments)
	require.Equal(t, err.Error(), "invalid arguments\nfoo bar is invalid")
}

func TestEBadRequest(t *testing.T) {
	err := EBadRequest("bar")
	require.ErrorIs(t, err, ErrInvalidFormat)
	require.Equal(t, "invalid format\nbar with invalid format", err.Error())
}
func TestHandleError(t *testing.T) {
	cases := []struct {
		Name               string
		Err                error
		ExpectedStatusCode int
	}{
		{
			Name:               "WHEN nil error",
			Err:                nil,
			ExpectedStatusCode: http.StatusInternalServerError,
		},
		{
			Name:               "WHEN ErrInvalidFormat",
			Err:                ErrInvalidFormat,
			ExpectedStatusCode: http.StatusBadRequest,
		},
		{
			Name:               "WHEN ErrInvalidArguments",
			Err:                ErrInvalidArguments,
			ExpectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			Name:               "WHEN ErrEmptyArguments",
			Err:                ErrEmptyArguments,
			ExpectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			Name:               "WHEN ErrConflict",
			Err:                ErrConflict,
			ExpectedStatusCode: http.StatusConflict,
		},
		{
			Name:               "WHEN no mapped error",
			Err:                ErrNotFound,
			ExpectedStatusCode: http.StatusNotFound,
		},
	}
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			w := httptest.NewRecorder()
			HandleError(w, c.Err)
			require.Equal(t, c.ExpectedStatusCode, w.Result().StatusCode)
		})
	}
}
