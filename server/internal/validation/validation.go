package validation

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net"
	"net/http"
	"regexp"
)

var (
	ErrInvalidDomain = errors.New("invalid domain")
	ErrUnknownDomain = errors.New("unknown domain")
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// BindAndValidate attempts to bind the form to the given struct and validates the result.
// Returns an HTTP error if binding or validation fails.
func BindAndValidate(c echo.Context, i interface{}) *echo.HTTPError {
	if err := c.Bind(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	if err := c.Validate(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func ValidateDomain(domain string) error {
	var domainRegex = regexp.MustCompile(`^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,}$`)
	if !domainRegex.MatchString(domain) {
		return ErrInvalidDomain
	}
	_, err := net.LookupHost(domain)
	if err != nil {
		return ErrUnknownDomain
	}
	return nil
}
