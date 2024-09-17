package sbext

import (
	"encoding/json"
	"strings"
)

// Error represents the JSON error returned from Supabase API calls.
type Error struct {
	Code      int    `json:"code"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"msg"`
}

// Error returns the message of the error.
func (e *Error) Error() string {
	return e.Message
}

// NewError creates a new *Error from an error.
// The method also returns a boolean indicating of the error was successfully deserialized.
func NewError(err error) (*Error, bool) {
	data := err.Error()
	start := strings.Index(data, "{")
	if start == -1 {
		return nil, false
	}
	jsonData := data[start:]

	var sbErr Error
	if err := json.Unmarshal([]byte(jsonData), &sbErr); err != nil {
		return nil, false
	}
	return &sbErr, sbErr.Code != 0
}

const (
	ErrorCodeInvalidCredentials string = "invalid_credentials"
	ErrorCodeOTPExpired         string = "otp_expired"
	ErrorCodeSamePassword       string = "same_password"
)
