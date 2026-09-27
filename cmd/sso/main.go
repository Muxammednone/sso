package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	ssov1 "github.com/Muxammednone/sso-proto/gen/go/sso"
	"github.com/Muxammednone/sso/internal/config"
	"github.com/Muxammednone/sso/internal/repository/postgres"
	"github.com/Muxammednone/sso/internal/service"
	"github.com/Muxammednone/sso/internal/token"
	transport_grpc "github.com/Muxammednone/sso/internal/transport/grpc"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

func main() {
	shutdownCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	config, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	connectCtx, cancel := context.WithTimeout(context.Background(), config.DBTimeout)
	defer cancel()

	pool, err := pgxpool.New(connectCtx, config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()

	if err := pool.Ping(connectCtx); err != nil {
		log.Printf("connect to PostgreSQL: %v", err)
		return
	}
	cancel()

	log.Println("connected to PostgreSQL")

	userRepository := postgres.NewUserRepository(pool)
	jwtToken := token.NewJWTGenerator([]byte(config.JWTSecret), config.TokenTTL)
	userService := service.NewUserService(userRepository, jwtToken)

	handler := transport_grpc.NewServer(userService)
	grpcServer := grpc.NewServer()
	ssov1.RegisterSSOServer(grpcServer, handler)

	listener, err := net.Listen("tcp", config.GRPCAddr)

	if err != nil {
		log.Printf("Listen tcp server on port '%v': %v", config.GRPCAddr, err)
		return
	}
	defer listener.Close()

	log.Printf("Listen tcp server on port: '%v'", config.GRPCAddr)

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- grpcServer.Serve(listener)
	}()

	select {
	case err := <-serveErr:
		grpcServer.Stop()
		if err != nil {
			log.Printf("serve: %v", err)
		}
		return
	case <-shutdownCtx.Done():
		log.Println("stopping gRPC server")
	}

	stopped := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(stopped)
	}()

	timer := time.NewTimer(config.ShutdownTimeout)
	defer timer.Stop()
	select {
	case <-stopped:
	case <-timer.C:
		log.Println("shutdown timeout: forcing gRPC server to stop")
		grpcServer.Stop()
		<-stopped
	}

	log.Println("gRPC server stopped")
}
