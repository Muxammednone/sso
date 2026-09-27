package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	ssov1 "github.com/Muxammednone/sso-proto/gen/go/sso"
	"github.com/Muxammednone/sso/internal/repository/postgres"
	"github.com/Muxammednone/sso/internal/service"
	"github.com/Muxammednone/sso/internal/token"
	transport_grpc "github.com/Muxammednone/sso/internal/transport/grpc"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	if err := pool.Ping(context.Background()); err != nil {
		log.Printf("connect to PostgreSQL: %v", err)
		return
	}

	log.Println("connected to PostgreSQL")

	userRepository := postgres.NewUserRepository(pool)
	jwtToken := token.NewJWTGenerator([]byte(os.Getenv("JWT_SECRET")), time.Hour)
	userService := service.NewUserService(userRepository, jwtToken)

	handler := transport_grpc.NewServer(userService)
	grpcServer := grpc.NewServer()
	ssov1.RegisterSSOServer(grpcServer, handler)

	listener, err := net.Listen("tcp", "127.0.0.1:50051")

	if err != nil {
		log.Printf("Listen tcp server on port '127.0.0.1:50051': %v", err)
		return
	}
	defer listener.Close()

	log.Print("Listen tcp server on port '127.0.0.1:50051'")

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Printf("serve 'listener': %v", err)
		return
	}
}
