package transport_grpc

import (
	"context"

	ssov1 "github.com/Muxammednone/sso-proto/gen/go/sso"
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
		return nil, err
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
		return nil, err
	}

	return &ssov1.LoginResponse{
		AccessToken: token,
	}, nil

}
