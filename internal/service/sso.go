package service

import (
	"context"
	"errors"

	"github.com/Muxammednone/sso/internal/repository"
)

var errNotImplemented = errors.New("not implemented")

type Service struct {
	userRepo repository.UserRepository
}

func New(userRepo repository.UserRepository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) Register(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {
	return 0, errNotImplemented
}

func (s *Service) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	return "", errNotImplemented
}
