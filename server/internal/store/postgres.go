package store

import (
	"errors"
	"fmt"
	"github.com/lib/pq"
)

// https://www.postgresql.org/docs/current/errcodes-appendix.html#ERRCODES-TABLE

type PgErr string

const (
	PgErrNone                PgErr = "-"
	PgErrCodeUniqueViolation PgErr = "unique_violation"
)

// Error implements the built-in error interface
func (pge PgErr) Error() string {
	return string(pge)
}

// String returns a string representation of the PgErr.
// Additional information is included if the PgErr is not a PostgresSQL
// error code or is unhandled.
func (pge PgErr) String() string {
	s := string(pge)
	if pge.Error() == string(PgErrNone) || s == "" {
		if s == "" {
			s = "empty-string"
		}
		return fmt.Sprintf("non-PostgresSQL error code: %s", s)
	}
	if _, ok := stringToPgErr(s); !ok {
		return fmt.Sprintf("unhandled PostgresSQL error code: %s", s)
	}
	return s
}

// stringToPgErr returns a PgErr from a given string code and a boolean
// indicating if the error code has been handled. A false response indicates
// there was a valid error code, but it is not handled. True indicates that
// either there is no error or the error code is handled.
func stringToPgErr(code string) (PgErr, bool) {
	switch code {
	case string(PgErrNone), "":
		return PgErrNone, false
	case "unique_violation":
		return PgErrCodeUniqueViolation, true
	default:
		return PgErr(code), false
	}
}

func checkPgErr(err error) PgErr {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return PgErrNone
	}
	pge, _ := stringToPgErr(pqErr.Code.Name())
	return pge
}
