package contract

import (
	"advancely/internal/model"
	"github.com/google/uuid"
)

type Store interface {
	UserStore
	CompanyStore
}

type UserStore interface {
	// User returns the user associated with the given id.
	User(id uuid.UUID) (model.UserProfile, error)
	// BaseUserByEmail returns the auth.users user associated with the given email.
	BaseUserByEmail(email string) (model.User, error)
	// Users returns a slice of all users.
	Users(companyID uuid.UUID) ([]model.UserProfile, error)
	// CreateProfile creates a record in the profiles table.
	CreateProfile(user *model.UserProfile) error
	UpdateUser(user *model.UserProfile) error
	DeleteUser(id uuid.UUID) error
}

type CompanyStore interface {
	// Company returns the company associated with the given id.
	Company(id uuid.UUID) (model.Company, error)
	// CompanyByCreator returns the company created by the given creator user ID.
	CompanyByCreator(creatorID uuid.UUID) (model.Company, error)
	// Companies returns a slice of all companies.
	Companies() ([]model.Company, error)
	CreateCompany(c *model.Company) error
	UpdateCompany(c *model.Company) error
	DeleteCompany(id uuid.UUID) error
}
