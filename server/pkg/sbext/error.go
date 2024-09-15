package sbext

import (
	"encoding/json"
	"io"
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

// NewError creates a new *Error from an io.Reader interface.
// The method also returns a boolean indicating of the error was successfully
// deserialized from the reader.
func NewError(r io.Reader) (*Error, bool) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, false
	}
	var sbErr Error
	if err = json.Unmarshal(b, &sbErr); err != nil {
		return nil, false
	}
	return &sbErr, true
}

const (
	SupabaseErrorCodeOTPExpired string = "otp_expired"
)
