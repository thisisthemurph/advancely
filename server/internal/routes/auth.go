package routes

import (
	"advancely/internal/auth"
	"advancely/internal/model"
	"advancely/internal/store"
	"advancely/internal/store/contract"
	"advancely/internal/validation"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/nedpals/supabase-go"
	"net/http"
)

func NewAuthHandler(sb *supabase.Client, s *store.Store, sessionSecret string) AuthHandler {
	return AuthHandler{
		Supabase:      sb,
		UserStore:     s.UserStore,
		CompanyStore:  s.CompanyStore,
		SessionSecret: sessionSecret,
	}
}

type AuthHandler struct {
	Supabase      *supabase.Client
	UserStore     contract.UserStore
	CompanyStore  contract.CompanyStore
	SessionSecret string
}

func (h AuthHandler) MakeRoutes(e *echo.Group) {
	group := e.Group("/auth")
	group.POST("/login", h.handleLogin())
	group.POST("/signup", h.handleSignup())
	group.POST("/confirm-email", h.handleVerifyEmailVerificationComplete())
}

func (h AuthHandler) handleLogin() echo.HandlerFunc {
	type Request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(c echo.Context) error {
		ctx := c.Request().Context()

		user := auth.Session(c)
		fmt.Println(user)

		var req Request
		if err := validation.BindAndValidate(c, &req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		authDetails, err := h.Supabase.Auth.SignIn(ctx, supabase.UserCredentials{
			Email:    req.Email,
			Password: req.Password,
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		session := auth.NewSessionCookie(authDetails)
		if err := session.SetCookie(c, h.SessionSecret); err != nil {
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
			return err
		}

		user, err := getOrSignupSupabaseUser(ctx, form)
		if err != nil {
			we := fmt.Errorf("error signing up user: %w", err)
			return echo.NewHTTPError(500, we.Error())
		}
		userID, err := uuid.Parse(user.ID)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Invalid user ID")
		}

		// CreateCompany the company

		company, err := getOrCreateCompany(model.Company{
			Name:      form.CompanyName,
			CreatorID: userID,
		})
		if err != nil {
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
			return err
		}

		user, err := h.Supabase.Auth.User(ctx, req.Token)
		if err != nil || user == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		if user.ConfirmedAt.IsZero() {
			return c.NoContent(http.StatusUnauthorized)
		}

		return c.NoContent(http.StatusNoContent)
	}
}
