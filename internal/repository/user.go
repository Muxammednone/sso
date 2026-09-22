package repository

import (
	"context"
	"errors"

	"github.com/Muxammednone/sso/internal/domain"
)

var (
	ErrNotFound = errors.New("user not found")
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
