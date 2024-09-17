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

const recoverPath string = "/recover"

func NewSupabaseExtended(client *supabase.Client, config application.SupabaseConfig, webBaseURL string) *SupabaseExtended {
	return &SupabaseExtended{
		Client: client,
		Extensions: &Extensions{
			Client:     client,
			Config:     config,
			WebBaseURL: webBaseURL,
		},
	}
}

// Extensions comprises a set of extended Supabase utility functions where
// the supabase-go project has missing or incomplete functionality.
type Extensions struct {
	*supabase.Client
	Config     application.SupabaseConfig
	WebBaseURL string
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

	redirectTo = fmt.Sprintf("%s%s", c.WebBaseURL, redirectTo)
	path := recoverPath + "?redirect_to=" + redirectTo

	resp, err := c.post(ctx, path, bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("supabase password recovery responded with unexpected status code: %d => %s", resp.StatusCode, string(b))
	}
	return nil
}

// post sends a POST request to the Supabase server and returns the response.
// The apikey header is set automatically on all requests.
func (c *Extensions) post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/auth/v1/%s", c.Config.PublicKey, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", c.Config.PublicKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to post to %s: %w", url, err)
	}
	return resp, nil
}
