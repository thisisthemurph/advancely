package sbext

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"advancely/internal/application"
	"github.com/supabase-community/supabase-go"
)

const (
	authPath    string = "/auth/v1"
	recoverPath string = "/recover"
)

func NewSupabaseExtended(client *supabase.Client, config application.SupabaseConfig) *SupabaseExtended {
	return &SupabaseExtended{
		Client: client,
		Extensions: &Extensions{
			Client:     client,
			apiKey:     config.PublicKey,
			baseURL:    config.URL,
			httpClient: http.DefaultClient,
		},
	}
}

// Extensions comprises a set of extended Supabase utility functions where
// the supabase-go project has missing or incomplete functionality.
type Extensions struct {
	*supabase.Client
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// SupabaseExtended wraps the supabase.Client, providing extended functionality.
type SupabaseExtended struct {
	*supabase.Client
	Extensions *Extensions
}

// ResetPasswordForEmail sends a password recovery link to the given e-mail address.
func (c *Extensions) ResetPasswordForEmail(ctx context.Context, email, redirectTo string) error {
	b, err := json.Marshal(map[string]string{
		"email": email,
	})
	if err != nil {
		return err
	}

	path := recoverPath
	if redirectTo != "" {
		path += "?redirect_to=" + redirectTo
	}
	r, err := c.newRequest(ctx, authPath+path, http.MethodPost, bytes.NewReader(b))
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(r)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if !statusCodeIsSuccess(resp.StatusCode) {
		fullBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("response status code: %d", resp.StatusCode)
		}
		return fmt.Errorf("response status code %d: %s", resp.StatusCode, string(fullBody))
	}
	return nil
}

func (c *Extensions) newRequest(ctx context.Context, path string, method string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("apikey", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func statusCodeIsSuccess(code int) bool {
	return code >= 200 && code < 300
}
