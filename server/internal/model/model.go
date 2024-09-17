package model

import (
	"github.com/google/uuid"
	"github.com/supabase-community/gotrue-go/types"
	"time"
)

// User represents a subset of columns from the Supabase auth.users table.
type User struct {
	ID                 uuid.UUID  `db:"id"`
	Aud                string     `db:"aud"`
	Role               string     `db:"role"`
	Email              string     `db:"email"`
	EmailConfirmedAt   *time.Time `db:"email_confirmed_at"`
	InvitedAt          *time.Time `db:"invited_at"`
	ConfirmationSentAt *time.Time `db:"confirmation_sent_at"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
}

// SupabaseUser converts a User to a gotrue-go types.User.
func (u *User) SupabaseUser() *types.User {
	return &types.User{
		ID:                 u.ID,
		Aud:                u.Aud,
		Role:               u.Role,
		Email:              u.Email,
		EmailConfirmedAt:   u.EmailConfirmedAt,
		InvitedAt:          u.InvitedAt,
		ConfirmationSentAt: u.ConfirmationSentAt,
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
	}
}

// UserProfile represents a combination columns from the auth.users and public.profile tables
type UserProfile struct {
	ID        uuid.UUID  `db:"id"`
	CompanyID uuid.UUID  `db:"company_id"`
	FirstName string     `db:"first_name"`
	LastName  string     `db:"last_name"`
	Email     string     `db:"email"`
	IsAdmin   bool       `db:"is_admin"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

type Company struct {
	ID        uuid.UUID  `db:"id"`
	Name      string     `db:"name"`
	CreatorID uuid.UUID  `db:"creator_id"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}
