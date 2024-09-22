package routes

import (
	"log/slog"
	"net/http"

	"advancely/internal/auth"
	"advancely/internal/store"
	"advancely/internal/validation"

	"github.com/labstack/echo/v4"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

func NewUsersHandler(
	s *store.PostgresStore,
	sb *supabase.Client,
	ensurePermissionFn EnsurePermissionFn,
	logger *slog.Logger) UsersHandler {
	return UsersHandler{
		UserStore:        s.UserStore,
		Supabase:         sb,
		EnsurePermission: ensurePermissionFn,
		Logger:           logger,
	}
}

type UsersHandler struct {
	UserStore        store.UserStore
	EnsurePermission EnsurePermissionFn
	Supabase         *supabase.Client
	Logger           *slog.Logger
}

func (h UsersHandler) MakeRoutes(e *echo.Group) {
	group := e.Group("/user")
	group.GET("", h.HandleListUsers())
	group.POST("", h.HandleCreateNewUser())
}

type NewUserRequest struct {
	FirstName string `json:"firstName" validation:"required"`
	LastName  string `json:"lastName" validation:"required"`
	Email     string `json:"email" validation:"required,email"`
}

func (h UsersHandler) HandleListUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		session := auth.CurrentUser(c)
		users, err := h.UserStore.Users(session.Company.ID)
		if err != nil {
			h.Logger.Error("error listing users", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		res := NewPagedResponse(c, users)
		return c.JSON(http.StatusOK, res)
	}
}

// HandleCreateNewUser adds a user to the company and sends an invitation email.
// A record is also added to the public.profiles table for the user.
func (h UsersHandler) HandleCreateNewUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		session := auth.CurrentUser(c)

		var req NewUserRequest
		if err := validation.BindAndValidate(c, &req); err != nil {
			return err
		}

		exists, err := h.UserStore.Exists(req.Email)
		if err != nil {
			h.Logger.Error("error checking if user already exists", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		if exists {
			return echo.NewHTTPError(http.StatusConflict, "A user with the given email already exists")
		}

		err = h.Supabase.Auth.OTP(types.OTPRequest{
			Email:      req.Email,
			CreateUser: true,
			Data: map[string]interface{}{
				"created_by": map[string]string{
					"id":    session.User.ID.String(),
					"email": session.User.Email,
				},
			},
		})

		if err != nil {
			h.Logger.Error("error creating user", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		user, err := h.UserStore.BaseUserByEmail(req.Email)
		if err != nil {
			h.Logger.Error("error fetching recently created user", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		profile, err := h.UserStore.CreateProfile(store.CreateProfileRequest{
			UserID:    user.ID,
			CompanyID: session.Company.ID,
			FirstName: req.FirstName,
			LastName:  req.LastName,
		})
		profile.Email = req.Email
		return c.JSON(http.StatusCreated, profile)
	}
}
