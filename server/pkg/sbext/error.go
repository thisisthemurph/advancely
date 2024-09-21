package sbext

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Error represents the JSON error returned from Supabase API calls.
type Error struct {
	Code      int    `json:"code"`
	ErrorCode string `json:"error"`
	Message   string `json:"error_description"`
}

// Error returns the message of the error.
func (e *Error) Error() string {
	return e.Message
}

// NewError creates a new *Error from an error.
// The method also returns a boolean indicating of the error was successfully deserialized.
func NewError(checkErr error) (*Error, bool) {
	data := checkErr.Error()
	re := regexp.MustCompile(`response status code (\d+):`)
	matches := re.FindStringSubmatch(data)
	if len(matches) < 2 {
		return nil, false
	}

	var statusCode int
	if _, err := fmt.Sscanf(matches[1], "%d", &statusCode); err != nil {
		return nil, false
	}

	start := strings.Index(data, "{")
	if start == -1 {
		return nil, false
	}
	jsonData := data[start:]

	var sbErr Error
	if err := json.Unmarshal([]byte(jsonData), &sbErr); err != nil {
		return nil, false
	}
	sbErr.Code = statusCode
	return &sbErr, sbErr.Code != 0
}

const (
	ErrorCodeInvalidCredentials string = "invalid_credentials"
	ErrorCodeOTPExpired         string = "otp_expired"
	ErrorCodeSamePassword       string = "same_password"
)
