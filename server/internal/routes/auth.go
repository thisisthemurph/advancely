package routes

import (
	"advancely/internal/application"
	"advancely/internal/auth"
	"advancely/internal/model"
	"advancely/internal/store"
	"advancely/internal/store/contract"
	"advancely/internal/validation"
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nedpals/supabase-go"
	"log/slog"
	"net/http"
)

func NewAuthHandler(
	sb *supabase.Client,
	s *store.Store,
	config application.AppConfig,
	logger *slog.Logger,
) AuthHandler {
	return AuthHandler{
		Supabase:     sb,
		UserStore:    s.UserStore,
		CompanyStore: s.CompanyStore,
		Config:       config,
		Logger:       logger,
	}
}

type AuthHandler struct {
	Supabase     *supabase.Client
	UserStore    contract.UserStore
	CompanyStore contract.CompanyStore
	Config       application.AppConfig
	Logger       *slog.Logger
}

func (h AuthHandler) MakeRoutes(e *echo.Group) {
	group := e.Group("/auth")
	group.POST("/login", h.handleLogin())
	group.POST("/signup", h.handleSignup())
	group.POST("/logout", h.handleLogout())
	group.POST("/confirm-email", h.handleVerifyEmailVerificationComplete())
}

func (h AuthHandler) handleLogout() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		user := auth.Session(c)
		if err := h.Supabase.Auth.SignOut(ctx, user.AccessToken); err != nil {
			h.Logger.Error("Error logging out", "error", err)
		}
		if err := auth.DeleteSessionCookie(c, h.Config.SessionSecret); err != nil {
			h.Logger.Error("Error deleting session cookie", "error", err)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func (h AuthHandler) handleLogin() echo.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(c echo.Context) error {
		ctx := c.Request().Context()

		var req Request
		if err := validation.BindAndValidate(c, &req); err != nil {
			h.Logger.Error("failed binding/validating login request", "error", err)
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		authDetails, err := h.Supabase.Auth.SignIn(ctx, supabase.UserCredentials{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			h.Logger.Error("failed signing in with supabase", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		session := auth.NewSessionCookie(authDetails)

		user, err := h.UserStore.User(session.Sub)
		if err != nil {
			h.Logger.Error("failed getting user from store", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		session.SetUser(user)

		company, err := h.CompanyStore.Company(user.CompanyID)
		if err != nil {
			h.Logger.Error("failed getting company from store", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		session.SetUserCompany(company)

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

	type UserMetadata struct {
		CompanyName string `json:"company_name"`
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
	getOrSignupSupabaseUser := func(ctx context.Context, form formParams) (*supabase.User, error) {
		existingUser, err := h.UserStore.BaseUserByEmail(form.UserEmail)
		if errors.Is(err, store.ErrUserNotFound) {
			return h.Supabase.Auth.SignUp(ctx, supabase.UserCredentials{
				Email:    form.UserEmail,
				Password: form.Password,
				Data:     UserMetadata{CompanyName: form.CompanyName},
			})
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
			return echo.NewHTTPError(500)
		}
		userID, err := uuid.Parse(user.ID)
		if err != nil {
			h.Logger.Error("failed to parse user ID", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user ID")
		}

		// CreateCompany the company

		company, err := getOrCreateCompany(model.Company{
			Name:      form.CompanyName,
			CreatorID: userID,
		})
		if err != nil {
			h.Logger.Error("failed to create company", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create company")
		}

		// CreateCompany the initial admin user profile

		profile, err := getOrCreateUserProfile(model.UserProfile{
			ID:        userID,
			CompanyID: company.ID,
			FirstName: form.UserFirstName,
			LastName:  form.UserLastName,
			Email:     user.Email,
			IsAdmin:   true,
		})
		if err != nil {
			h.Logger.Error("failed to create user profile", "error", err)
			return echo.NewHTTPError(500, "Failed to create profile")
		}

		return c.JSON(http.StatusOK, SignupResponse{
			ID:            user.ID,
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
		ctx := c.Request().Context()

		var req Request
		if err := validation.BindAndValidate(c, &req); err != nil {
			h.Logger.Error("failed to bind/validate verification request", "error", err)
			return err
		}

		user, err := h.Supabase.Auth.User(ctx, req.Token)
		if err != nil || user == nil {
			h.Logger.Error("failed to verify user", "error", err)
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		if user.ConfirmedAt.IsZero() {
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusNoContent)
	}
}
