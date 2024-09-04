package types

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	IsActive  bool      `json:"is_active"`
	UserType  string    `json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserPayload struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email" validate:"required"`
	Password  string    `json:"password" validate:"required"`
	IsActive  bool      `json:"is_active"`
	UserType  string    `json:"user_type" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PasswordUserPayload struct {
	Password  string    `json:"password" validate:"required"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterUserPayload struct {
	// FirstName string `json:"firstName" validate:"required"`
	// LastName  string `json:"lastName" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3,max=130"`
}

type CreateLoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Email     string    `json:"email"`
	IsActive  bool      `json:"is_active"`
	UserType  string    `json:"user_type"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type DisableUserPayload struct {
	ID        uuid.UUID `json:"id"`
	IsActive  bool      `json:"is_active"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserStore interface {
	CreateUser(UserPayload) error
	GetUsers() ([]*User, error)
	GetUserByID(id uuid.UUID) (*User, error)
	GetUserByEmail(email string) (*User, error)
	LoginUser(user CreateLoginPayload) (*User, error)
	// DisableUser(ctx context.Context, user UserPayload) error
	DisableUser(ctx context.Context, user DisableUserPayload) error
	// UpdateUser(User) error
}
