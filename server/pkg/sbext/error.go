package sbext

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	ErrorCodeInvalidCredentials string = "invalid_credentials"
	ErrorCodeOTPExpired         string = "otp_expired"
	ErrorCodeSamePassword       string = "same_password"
)

type errorJSON1 struct {
	Code      int    `json:"code"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"msg"`
}

type errorJSON2 struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

// Error represents the JSON error returned from Supabase API calls.
type Error struct {
	Code      int
	ErrorCode string
	Message   string
}

var parsers = []func(error) (*Error, error){parseErrorJSON1, parseErrorJSON2}

// NewError creates a new *Error from an error.
// The method also returns a boolean indicating of the error was successfully deserialized.
func NewError(errString error) (*Error, bool) {
	for _, parser := range parsers {
		if res, err := parser(errString); err == nil {
			return res, true
		}
	}
	return nil, false
}

// Error returns the message of the error.
func (e *Error) Error() string {
	return e.Message
}

// parseErrorJSON1 parses an *Error from a string error that has the structure:
//
//	response status code 400: {"code":400,"error_code":"str_err_code","msg":"Human error message"}
func parseErrorJSON1(e error) (*Error, error) {
	errString := e.Error()
	if !strings.HasPrefix(errString, "response status code") {
		return nil, errors.New("unrecognised error")
	}
	if !strings.Contains(errString, `"code"`) &&
		!strings.Contains(errString, `"error_code"`) &&
		!strings.Contains(errString, `"msg"`) {
		return nil, errors.New("unrecognised error")
	}

	start := strings.Index(errString, "{")
	if start == -1 {
		return nil, errors.New("not a JSON string")
	}
	errJSON := errString[start:]

	var res errorJSON1
	if err := json.Unmarshal([]byte(errJSON), &res); err != nil {
		return nil, err
	}
	return &Error{
		Code:      res.Code,
		ErrorCode: res.ErrorCode,
		Message:   res.Message,
	}, nil
}

// parseErrorJSON2 parses an *Error from a string error that has the structure:
//
//	response status code 400: {"error":"str_error_code","error_description":"Human error message"}
func parseErrorJSON2(e error) (*Error, error) {
	errString := e.Error()
	if !strings.HasPrefix(errString, "response status code") {
		return nil, errors.New("unrecognised error")
	}
	if !strings.Contains(errString, `"error"`) && !strings.Contains(errString, `"error_description"`) {
		return nil, errors.New("unrecognised error")
	}

	statusCode, err := parseStatusCode(errString)
	if err != nil {
		return nil, err
	}

	start := strings.Index(errString, "{")
	if start == -1 {
		return nil, errors.New("not a JSON string")
	}
	errJSON := errString[start:]

	var res errorJSON2
	if err := json.Unmarshal([]byte(errJSON), &res); err != nil {
		return nil, err
	}
	return &Error{
		Code:      statusCode,
		ErrorCode: res.Error,
		Message:   res.ErrorDescription,
	}, nil
}

func parseStatusCode(s string) (int, error) {
	re := regexp.MustCompile(`response status code (\d+):`)
	matches := re.FindStringSubmatch(s)
	if len(matches) < 2 {
		return 0, errors.New("bad string")
	}
	statusCode, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, fmt.Errorf("bad string: %w", err)
	}
	return statusCode, nil
}
