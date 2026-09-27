package service

import (
	"context"
	"errors"

	"github.com/Muxammednone/sso/internal/domain"
	"github.com/Muxammednone/sso/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailRequired      = errors.New("email is required")
	ErrPasswordRequired   = errors.New("password is required")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type Service struct {
	userRepository repository.UserRepository
	tokenGenerator TokenGenerator
}

func NewUserService(
	userRepository repository.UserRepository,
	tokenGenerator TokenGenerator,
) *Service {
	return &Service{
		userRepository: userRepository,
		tokenGenerator: tokenGenerator,
	}
}

type TokenGenerator interface {
	Generate(
		userID int64,
	) (string, error)
}

func (s *Service) Register(
	ctx context.Context,
	email string,
	password string,
) (int64, error) {

	if email == "" {
		return 0, ErrEmailRequired
	}

	if password == "" {
		return 0, ErrPasswordRequired
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := &domain.User{
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	createdUser, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return 0, err
	}

	return createdUser.ID, nil
}

func (s *Service) Login(
	ctx context.Context,
	email string,
	password string,
) (string, error) {
	// *Валидация
	if email == "" {
		return "", ErrEmailRequired
	}

	if password == "" {
		return "", ErrPasswordRequired
	}

	user, err := s.userRepository.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", ErrInvalidCredentials
	}

	return s.tokenGenerator.Generate(user.ID)
}
