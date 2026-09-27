package transport_grpc

import (
	"context"
	"errors"
	"log"

	ssov1 "github.com/Muxammednone/sso-proto/gen/go/sso"
	"github.com/Muxammednone/sso/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SSOService interface {
	Register(
		ctx context.Context,
		email string,
		password string,
	) (int64, error)

	Login(
		ctx context.Context,
		email string,
		password string,
	) (string, error)
}

type Server struct {
	ssov1.UnimplementedSSOServer

	service SSOService
}

func NewServer(service SSOService) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {

	userID, err := s.service.Register(
		ctx,
		req.GetEmail(),
		req.GetPassword(),
	)

	if err != nil {
		if errors.Is(err, service.ErrEmailRequired) || errors.Is(err, service.ErrPasswordRequired) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		log.Printf("register failed: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *Server) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	token, err := s.service.Login(
		ctx,
		req.GetEmail(),
		req.GetPassword(),
	)

	if err != nil {
		if errors.Is(err, service.ErrEmailRequired) || errors.Is(err, service.ErrPasswordRequired) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.Unauthenticated, err.Error())
		}
		log.Printf("login failed: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &ssov1.LoginResponse{
		AccessToken: token,
	}, nil

}
