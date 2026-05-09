package main

import (
	"fmt"
	"log"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	pb "github.com/pbuser/genproto/user"
	"github.com/pbuser/server/middleware"
	"github.com/pbuser/server/service"
	"google.golang.org/grpc"
)

const (
	Addr           = ":8080"
	NetWork        = "tcp"
	DefaultTimeout = 5 * time.Second
)

func main() {
	lister, err := net.Listen(NetWork, Addr)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
		return
	}

	defer lister.Close()
	defer middleware.CloseLogger()

	fmt.Println("server lister is ", lister.Addr())

	//拦截器，可注册日志，授权认证
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			middleware.TimeoutStreamInterceptor(DefaultTimeout),
			grpc_auth.StreamServerInterceptor(middleware.AuthInterceptor),
			grpc_zap.StreamServerInterceptor(middleware.ZapInterceptor()),
		)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			middleware.TimeoutUnaryInterceptor(DefaultTimeout),
			grpc_auth.UnaryServerInterceptor(middleware.AuthInterceptor),
			grpc_zap.UnaryServerInterceptor(middleware.ZapInterceptor()),
		)),
	)

	pb.RegisterUserServiceServer(grpcServer, service.NewUserService())
	pb.RegisterStreamServiceServer(grpcServer, service.NewStreamService())
	pb.RegisterStreamClientServer(grpcServer, service.NewUploadService())
	pb.RegisterStreamServer(grpcServer, service.NewBothStreamServer())

	pb.RegisterGoodServer(grpcServer, service.NewGoodService())

	err = grpcServer.Serve(lister)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}
