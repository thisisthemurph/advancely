package sbext

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/supabase-community/gotrue-go/types"
)

// ParseSessionFromErrJson is a helper function for parsing a gotrue-go types.Session
// from an error. For some reason, the community-supabase VerifyForUser method returns
// a JSON response in the returned error.
func ParseSessionFromErrJson(err error) (*types.Session, error) {
	data := err.Error()

	if !strings.HasPrefix(data, "response status code 200") {
		return nil, errors.New("not a success response")
	}
	start := strings.Index(data, "{")
	if start == -1 {
		return nil, errors.New("data is not valid JSON")
	}

	jsonData := data[start:]
	var session types.Session
	if err := json.Unmarshal([]byte(jsonData), &session); err != nil {
		return nil, errors.New("error parsing session JSON")
	}
	return &session, nil
}
