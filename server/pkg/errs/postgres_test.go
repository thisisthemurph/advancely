package errs

import (
	"errors"
	"github.com/lib/pq"
	"testing"
)

func TestCheckPgErr(t *testing.T) {
	tests := []struct {
		err      error
		expected PgErr
	}{
		{&pq.Error{Code: "23505"}, PgErrCodeUniqueViolation},
		{errors.New("some other error"), PgErrNone},
		{nil, PgErrNone},                      // Check if nil returns PgErrNone
		{&pq.Error{Code: "99999"}, PgErrNone}, // Check an unhandled error code
	}

	for _, test := range tests {
		result := CheckPgErr(test.err)
		if result != test.expected {
			t.Errorf("CheckPgErr(%q) = %v; want %v", test.err, result, test.expected)
		}
	}
}

func TestStringToPgErr(t *testing.T) {
	tests := []struct {
		code     string
		expected PgErr
		handled  bool
	}{
		{"unique_violation", PgErrCodeUniqueViolation, true},
		{"", PgErrNone, false},
		{"99999", PgErr("99999"), false},
	}

	for _, test := range tests {
		result, handled := stringToPgErr(test.code)
		if result != test.expected || handled != test.handled {
			t.Errorf("stringToPgErr(%q) = %v, %v; want %v, %v", test.code, result, handled, test.expected, test.handled)
		}
	}
}

func TestPgErrString(t *testing.T) {
	tests := []struct {
		err      PgErr
		expected string
	}{
		{PgErr(""), "non-PostgresSQL error code: empty-string"},
		{PgErrNone, "non-PostgresSQL error code"},
		{PgErrCodeUniqueViolation, "unique_violation"},
		{PgErr("99999"), "unhandled PostgresSQL error code: 99999"},
	}

	for _, test := range tests {
		result := test.err.String()
		if result != test.expected {
			t.Errorf("PgErr(%q).String() = %q; want %q", test.err, result, test.expected)
		}
	}
}
