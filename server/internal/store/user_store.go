package store

import (
	"advancely/internal/model"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var ErrUserNotFound = errors.New("user not found")

func NewPostgresUserStore(db *sqlx.DB) *PostgresUserStore {
	return &PostgresUserStore{
		DB: db,
	}
}

type PostgresUserStore struct {
	*sqlx.DB
}

func (s *PostgresUserStore) User(id uuid.UUID) (model.UserProfile, error) {
	var u model.UserProfile
	query := `
		select 
		    u.id, p.company_id, p.first_name, p.last_name, 
		    u.email, p.is_admin, p.created_at, p.updated_at 
		from auth.users u
		join public.profiles p on u.id = p.id
		where u.id = $1
		limit 1;`

	if err := s.Get(&u, query, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.UserProfile{}, ErrUserNotFound
		}
		return model.UserProfile{}, fmt.Errorf("error getting user: %w", err)
	}
	return u, nil
}

func (s *PostgresUserStore) BaseUserByEmail(email string) (model.User, error) {
	query := `
		select id, aud, role, email, email_confirmed_at, invited_at,
		       confirmation_sent_at, created_at, updated_at
		from auth.users
		where email = $1 
		limit 1;`

	var u model.User
	if err := s.Get(&u, query, email); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, ErrUserNotFound
		}
		return model.User{}, fmt.Errorf("error getting user: %w", err)
	}

	return u, nil
}

func (s *PostgresUserStore) Exists(email string) (bool, error) {
	var exists bool
	stmt := "select exists(select 1 from auth.users where email = $1);"
	if err := s.Get(&exists, stmt, email); err != nil {
		return exists, err
	}
	return exists, nil
}

func (s *PostgresUserStore) Users(companyID uuid.UUID) ([]model.UserProfile, error) {
	var uu []model.UserProfile
	query := `
		select 
		    u.id, p.company_id, p.first_name, p.last_name, 
		    u.email, p.is_admin, p.created_at, p.updated_at 
		from auth.users u
		join public.profiles p on u.id = p.id
		where p.company_id = $1;`

	if err := s.Select(&uu, query, companyID); err != nil {
		return []model.UserProfile{}, err
	}
	return uu, nil
}

type CreateProfileRequest struct {
	UserID    uuid.UUID
	CompanyID uuid.UUID
	FirstName string
	LastName  string
	IsAdmin   bool
}

func (s *PostgresUserStore) CreateProfile(req CreateProfileRequest) (model.UserProfile, error) {
	query := `
		insert into public.profiles (id, company_id, first_name, last_name, is_admin)
		values ($1, $2, $3, $4, $5)
		returning id, company_id, first_name, last_name, is_admin, created_at, updated_at;`

	var profile model.UserProfile
	if err := s.Get(&profile, query, req.UserID, req.CompanyID, req.FirstName, req.LastName, req.IsAdmin); err != nil {
		// TODO: Check if the profile already exists
		return model.UserProfile{}, fmt.Errorf("error creating profile: %w", err)
	}
	return profile, nil
}

func (s *PostgresUserStore) UpdateUser(user *model.UserProfile) error {
	query := `
		update public.profiles 
		set first_name = $1, last_name = $2, is_admin = $3
		where id = $4
		returning *;`

	if err := s.Get(user, query, user.FirstName, user.LastName, user.IsAdmin, user.ID); err != nil {
		return fmt.Errorf("error updating profile: %w", err)
	}
	return nil
}

func (s *PostgresUserStore) DeleteUser(id uuid.UUID) error {
	if _, err := s.Exec("delete from auth.users where id = $1;", id); err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	return nil
}
