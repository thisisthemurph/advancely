package sbext

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewError(t *testing.T) {
	testCases := []struct {
		name  string
		error error
	}{
		{
			name:  "success: JSON error type 1",
			error: errors.New(`response status code 400: {"code":400,"error_code":"str_err_code","msg":"Human error message"}`),
		},
		{
			name:  "success: JSON error type 2",
			error: errors.New(`response status code 400: {"error":"str_err_code","error_description":"Human error message"}`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, ok := NewError(tc.error)
			require.True(t, ok)
			require.Equal(t, 400, res.Code)
			require.Equal(t, "Human error message", res.Message)
			require.Equal(t, "str_err_code", res.ErrorCode)
		})
	}
}

func TestNewErrorWithInvalidErrors(t *testing.T) {
	testCases := []struct {
		name  string
		error error
	}{
		{
			name:  "failure: empty error string",
			error: errors.New(""),
		},
		{
			name:  "failure: partial JSON",
			error: errors.New(`response status code 400: {"code":400`),
		},
		{
			name:  "failure: missing closing brace",
			error: errors.New(`response status code 400: {"code":400,"error_code":"str_err_code","msg":"Human error message"`),
		},
		{
			name:  "failure: different properties",
			error: errors.New(`response status code 400: {"code":400,"error_value":"str_err_code","message":"Human error message"`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, ok := NewError(tc.error)
			require.False(t, ok)
			require.Nil(t, res)
		})
	}
}
