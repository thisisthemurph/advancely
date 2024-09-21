package auth_test

import (
	"advancely/internal/auth"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestExpired(t *testing.T) {
	testCases := []struct {
		name              string
		expiresAtUnixTime int64
		expectedResult    bool
	}{
		{
			name:              "expired time",
			expiresAtUnixTime: time.Now().UTC().Add(-time.Hour * 1).Unix(),
			expectedResult:    true,
		},
		{
			name:              "non-expired time",
			expiresAtUnixTime: time.Now().UTC().Add(time.Hour * 1).Unix(),
			expectedResult:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			session := auth.SessionCookie{
				ExpiresAt: tc.expiresAtUnixTime,
			}

			require.Equal(t, tc.expectedResult, session.Expired())
		})
	}
}
