package routes

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"advancely/internal/application"
	"advancely/internal/auth"
	"advancely/internal/model"
	"advancely/internal/model/security"
	"advancely/internal/store"
	"advancely/internal/validation"
	"advancely/pkg/sbext"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

func NewAuthHandler(
	supabaseClient *supabase.Client,
	s *store.PostgresStore,
	config application.AppConfig,
	logger *slog.Logger,
) AuthHandler {
	return AuthHandler{
		Supabase:         sbext.NewSupabaseExtended(supabaseClient, config.Supabase),
		UserStore:        s.UserStore,
		CompanyStore:     s.CompanyStore,
		PermissionsStore: s.PermissionsStore,
		Config:           config,
		Logger:           logger,
	}
}

type AuthHandler struct {
	Supabase         *sbext.SupabaseExtended
	UserStore        store.UserStore
	CompanyStore     store.CompanyStore
	PermissionsStore store.PermissionsStore
	Config           application.AppConfig
	Logger           *slog.Logger
}

func (h AuthHandler) MakeRoutes(e *echo.Group) {
	group := e.Group("/auth")
	group.POST("/login", h.HandleLogin())
	group.POST("/signup", h.handleSignup())
	group.POST("/logout", h.handleLogout())
	group.POST("/confirm-email", h.handleVerifyEmailVerificationComplete())
	group.POST("/reset-password", h.HandleTriggerPasswordReset())
	group.POST("/reset-password/confirm", h.handleConfirmPasswordReset())
}

func (h AuthHandler) handleLogout() echo.HandlerFunc {
	return func(c echo.Context) error {
		user := auth.CurrentUser(c)
		if err := h.Supabase.Auth.WithToken(user.AccessToken).Logout(); err != nil {
			h.Logger.Error("Error logging out", "error", err)
		}
		if err := auth.DeleteSessionCookie(c, h.Config.SessionSecret); err != nil {
			h.Logger.Error("Error deleting session cookie", "error", err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (h AuthHandler) HandleLogin() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req LoginRequest
		if err := validation.BindAndValidate(c, &req); err != nil {
			h.Logger.Error("failed binding/validating login request", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		token, err := h.Supabase.Auth.SignInWithEmailPassword(req.Email, req.Password)
		if err != nil {
			if sbErr, isSbErr := sbext.NewError(err); isSbErr {
				return echo.NewHTTPError(sbErr.Code, sbErr.Message)
			}
			h.Logger.Error("failed signing in with supabase", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		session := auth.NewSessionCookie(token.Session)

		user, err := h.UserStore.User(token.User.ID)
		if err != nil {
			if errors.Is(err, store.ErrUserNotFound) {
				return echo.NewHTTPError(http.StatusBadRequest, store.ErrUserNotFound)
			}
			h.Logger.Error("failed getting user from store", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		company, err := h.CompanyStore.Company(user.CompanyID)
		if err != nil {
			if errors.Is(err, store.ErrCompanyNotFound) {
				return echo.NewHTTPError(http.StatusBadRequest, store.ErrCompanyNotFound)
			}
			h.Logger.Error("failed getting company from store", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		session.SetUser(user)
		session.SetCompany(company)

		if err := session.SetCookie(c, h.Config.SessionSecret, h.Config.Environment); err != nil {
			h.Logger.Error("failed setting session cookie", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.JSON(http.StatusOK, session)
	}
}

// handleSignup signs the user up via Supabase and adds records to the companies and profiles tables.
// This function handles instances where the auth.user, company and profile rows may already exist
// due to any previous errored runs.
func (h AuthHandler) handleSignup() echo.HandlerFunc {
	type formParams struct {
		CompanyName   string `json:"name" validate:"required"`
		UserFirstName string `json:"firstName" validate:"required"`
		UserLastName  string `json:"lastName" validate:"required"`
		UserEmail     string `json:"email" validate:"required,email"`
		Password      string `json:"password" validate:"required,min=8"`
	}

	type SignupResponse struct {
		ID            string `json:"id"`
		UserFirstName string `json:"firstName"`
		UserLastName  string `json:"lastName"`
		UserEmail     string `json:"email"`
		CompanyName   string `json:"companyName"`
	}

	// getOrSignupSupabaseUser is an idempotent function for signing up with Supabase auth.
	// If the user already exists, the supabase.User is returned, otherwise signup is completed.
	getOrSignupSupabaseUser := func(ctx context.Context, form formParams) (*types.User, error) {
		existingUser, err := h.UserStore.BaseUserByEmail(form.UserEmail)
		if errors.Is(err, store.ErrUserNotFound) {
			resp, err := h.Supabase.Auth.Signup(types.SignupRequest{
				Email:    form.UserEmail,
				Password: form.Password,
				Data: map[string]interface{}{
					"company_name": form.CompanyName,
				},
			})
			// Signup returns either different response depending on the autoconfirm setting.
			// User if autoconfirm is off, Session if on.
			return &resp.User, err
		}
		return existingUser.SupabaseUser(), err
	}

	getOrCreateCompany := func(company model.Company) (model.Company, error) {
		existingCompany, err := h.CompanyStore.CompanyByCreator(company.CreatorID)
		if err == nil {
			return existingCompany, nil
		}
		if !errors.Is(err, store.ErrCompanyNotFound) {
			return model.Company{}, err
		}
		if err := h.CompanyStore.CreateCompany(&company); err != nil {
			return model.Company{}, err
		}
		return company, nil
	}

	getOrCreateUserProfile := func(user model.UserProfile) (model.UserProfile, error) {
		existingUser, err := h.UserStore.User(user.ID)
		if err == nil {
			return existingUser, nil
		}
		if !errors.Is(err, store.ErrUserNotFound) {
			return model.UserProfile{}, err
		}
		if err := h.UserStore.CreateProfile(&user); err != nil {
			return model.UserProfile{}, err
		}
		return user, nil
	}

	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var form formParams
		if err := validation.BindAndValidate(c, &form); err != nil {
			h.Logger.Error("failed to bind/validate signup request", "error", err)
			return err
		}

		user, err := getOrSignupSupabaseUser(ctx, form)
		if err != nil {
			h.Logger.Error("failed to sign up in supabase", "error", err)
			msg := fmt.Sprintf("Failed to create user with email address %s.", form.UserEmail)
			return echo.NewHTTPError(500, msg)
		}
		if err != nil {
			h.Logger.Error("failed to parse user ID", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user ID.")
		}

		// CreateCompany the company

		company, err := getOrCreateCompany(model.Company{
			Name:      form.CompanyName,
			CreatorID: user.ID,
		})
		if err != nil {
			h.Logger.Error("failed to create company", "error", err)
			msg := fmt.Sprintf("Failed to create company with name %s.", form.CompanyName)
			return echo.NewHTTPError(http.StatusInternalServerError, msg)
		}

		// CreateCompany the initial admin user profile

		profile, err := getOrCreateUserProfile(model.UserProfile{
			ID:        user.ID,
			CompanyID: company.ID,
			FirstName: form.UserFirstName,
			LastName:  form.UserLastName,
			Email:     user.Email,
			IsAdmin:   true,
		})
		if err != nil {
			h.Logger.Error("failed to create user profile", "error", err)
			return echo.NewHTTPError(500, "Failed to create your user profile")
		}

		// Assign the Admin role to the user

		err = h.PermissionsStore.AssignSystemRoleToUser(security.RoleAdmin, user.ID, company.ID)
		if err != nil {
			h.Logger.Error("failed to add admin role to user", "error", err)
			return echo.NewHTTPError(500, "Failed to assign appropriate permissions to your user")
		}

		return c.JSON(http.StatusOK, SignupResponse{
			ID:            user.ID.String(),
			UserFirstName: profile.FirstName,
			UserLastName:  profile.LastName,
			UserEmail:     user.Email,
			CompanyName:   company.Name,
		})
	}
}

func (h AuthHandler) handleVerifyEmailVerificationComplete() echo.HandlerFunc {
	type Request struct {
		Token string `json:"token" validate:"required"`
	}

	return func(c echo.Context) error {
		var req Request
		if err := validation.BindAndValidate(c, &req); err != nil {
			h.Logger.Error("failed to bind/validate verification request", "error", err)
			return err
		}

		user, err := h.Supabase.Auth.WithToken(req.Token).GetUser()
		if err != nil || user == nil || user.ID == uuid.Nil {
			h.Logger.Error("failed to verify user", "error", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		if user.EmailConfirmedAt.IsZero() {
			return echo.NewHTTPError(http.StatusUnauthorized, "Email address not confirmed")
		}
		return c.NoContent(http.StatusNoContent)
	}
}

type TriggerPasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

func (h AuthHandler) HandleTriggerPasswordReset() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req TriggerPasswordResetRequest
		if err := validation.BindAndValidate(c, &req); err != nil {
			return err
		}

		redirect := h.Config.ClientBaseURL + "/auth/password-reset"
		if err := h.Supabase.Extensions.ResetPasswordForEmail(ctx, req.Email, redirect); err != nil {
			h.Logger.Error("failed to trigger password reset", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

type OneTimePassword struct {
	OTP string `json:"otp"`
}

type ConfirmPasswordResetRequest struct {
	OneTimePassword
	Email       string `json:"email"`
	NewPassword string `json:"password"`
}

func (h AuthHandler) handleConfirmPasswordReset() echo.HandlerFunc {
	return func(c echo.Context) error {
		var req ConfirmPasswordResetRequest
		if err := validation.BindAndValidate(c, &req); err != nil {
			return err
		}

		// redirectTo seems to serve no purpose here, but is required.
		_, err := h.Supabase.Auth.VerifyForUser(types.VerifyForUserRequest{
			Type:       "recovery",
			Token:      req.OTP,
			Email:      req.Email,
			RedirectTo: h.Config.ClientBaseURL + "/auth/password-reset",
		})

		// The response from the above Supabase method seems to come in the error as JSON.
		// First we attempt to parse a success response from the error.
		session, jsonParseErr := sbext.ParseSessionFromErrJson(err)
		if jsonParseErr != nil && err != nil {
			sbErr, isSbErr := sbext.NewError(err)
			if isSbErr && sbErr.ErrorCode == sbext.ErrorCodeOTPExpired {
				return echo.NewHTTPError(http.StatusUnauthorized, "OTP is invalid or has expired")
			}
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		if session == nil || session.AccessToken == "" {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		cookie := auth.NewSessionCookie(*session)
		if err := cookie.SetCookie(c, h.Config.SessionSecret, h.Config.Environment); err != nil {
			h.Logger.Error("failed to set session cookie", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		_, err = h.Supabase.Auth.WithToken(session.AccessToken).UpdateUser(types.UpdateUserRequest{
			Password: &req.NewPassword,
		})
		if err != nil {
			sbErr, isSbErr := sbext.NewError(err)
			if isSbErr && sbErr.ErrorCode == sbext.ErrorCodeSamePassword {
				return echo.NewHTTPError(http.StatusBadRequest, sbErr.Message)
			}
			h.Logger.Error("failed to update user password", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
