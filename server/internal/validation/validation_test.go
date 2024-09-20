package validation

import (
	"testing"
)

func TestValidateDomain(t *testing.T) {
	tests := []struct {
		domain      string
		expectedErr error
	}{
		{"google.com", nil},
		{"outlook.co.uk", nil},
		{"", ErrInvalidDomain},
		{".", ErrInvalidDomain},
		{".com", ErrInvalidDomain},
		{".co.uk", ErrInvalidDomain},
		{"myname@gmail.com", ErrInvalidDomain},
		{"rbisgviwbweiycgwebyicwegvbcy.poi", ErrUnknownDomain},
	}

	for _, test := range tests {
		err := ValidateDomain(test.domain)
		if err != test.expectedErr {
			t.Errorf("ValidateDomain(%q) = %v; want %v", test.domain, err, test.expectedErr)
		}
	}
}
