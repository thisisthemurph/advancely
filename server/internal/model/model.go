package model

import (
	"github.com/google/uuid"
	"github.com/nedpals/supabase-go"
	"time"
)

// User represents a subset of columns from the Supabase auth.users table.
type User struct {
	ID                 string     `db:"id"`
	Aud                string     `db:"aud"`
	Role               string     `db:"role"`
	Email              string     `db:"email"`
	InvitedAt          *time.Time `db:"invited_at"`
	ConfirmedAt        *time.Time `db:"confirmed_at"`
	ConfirmationSentAt *time.Time `db:"confirmation_sent_at"`
	CreatedAt          time.Time  `db:"created_at"`
	UpdatedAt          time.Time  `db:"updated_at"`
}

func (u *User) SupabaseUser() *supabase.User {
	var invitedAt, confirmedAt, confirmationSentAt time.Time
	if u.InvitedAt != nil {
		invitedAt = *u.InvitedAt
	}
	if u.ConfirmedAt != nil {
		confirmedAt = *u.ConfirmedAt
	}
	if u.ConfirmationSentAt != nil {
		confirmationSentAt = *u.ConfirmationSentAt
	}
	return &supabase.User{
		ID:                 u.ID,
		Aud:                u.Aud,
		Role:               u.Role,
		Email:              u.Email,
		InvitedAt:          invitedAt,
		ConfirmedAt:        confirmedAt,
		ConfirmationSentAt: confirmationSentAt,
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
