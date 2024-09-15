package sbext

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"advancely/internal/application"
	"github.com/nedpals/supabase-go"
)

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
	recoverPasswordURL := fmt.Sprintf("%s/auth/v1/recover?redirect_to=%s", c.BaseURL, redirectTo)
	resp, err := c.post(ctx, recoverPasswordURL, bytes.NewReader(b))
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

// VerifyOTPForEmail logs the user in via the provided OTP token.
func (c *Extensions) VerifyOTPForEmail(ctx context.Context, email, otp string) (*supabase.AuthenticatedDetails, error) {
	body, err := json.Marshal(map[string]string{
		"type":  "recovery",
		"email": email,
		"token": otp,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to verify OTP: %w", err)
	}

	verifyTokenURL := fmt.Sprintf("%s/auth/v1/verify", c.BaseURL)
	resp, err := c.post(ctx, verifyTokenURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		supabaseErr, isSbErr := NewError(resp.Body)
		if !isSbErr {
			return nil, fmt.Errorf("OTP verification responded with unexpected status code: %d => %s", resp.StatusCode, string(body))
		}
		return nil, supabaseErr
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, errors.New("received empty response body")
	}

	var authenticatedDetails *supabase.AuthenticatedDetails
	err = json.Unmarshal(body, &authenticatedDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to parse auth details: %w", err)
	}

	return authenticatedDetails, nil
}

// post sends a POST request to the Supabase server and returns the response.
// The apikey header is set automatically on all requests.
func (c *Extensions) post(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", c.Config.PublicKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to post to %s: %w", url, err)
	}
	return resp, nil
}
