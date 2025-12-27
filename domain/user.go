package domain

import (
	"context"
	"errors"
	"time"
)

type User struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
type UserCriteria struct {
	ID       *uint
	Email    *string
	Name     *string
	Password *string
}

type UserUsecase interface {
	Create(ctx context.Context, user *User) (*User, error)
	Get(ctx context.Context, ctr *UserCriteria) (*User, error)
}
type UserRepository interface {
	Create(ctx context.Context, user *User) (*User, error)
	Get(ctx context.Context, ctr *UserCriteria) (*User, error)
}

var (
	ErrUserNotFound = errors.New("user not found")
)
