package grpc_main

import (
	"net"

	userpb "github.com/mi6ule/goboy/gen/go/proto/user/v1"
	"github.com/mi6ule/goboy/internal/infrastructure/database/persistence"
	errorhandler "github.com/mi6ule/goboy/internal/infrastructure/error-handler"
	grpc_service "github.com/mi6ule/goboy/internal/infrastructure/grpc/service"
	"google.golang.org/grpc"
)

func StartGRPCServer(db *persistence.MongoDatabase, redisClient *persistence.RedisClient) {
	lis, err := net.Listen("tcp", "127.0.0.1:9879")
	errorhandler.ErrorHandler(errorhandler.ErrorInput{Err: err, Message: "failed to listen GRPC", ErrType: "Fatal"})

	grpcUserService := grpc_service.NewGrpcUserService(db, redisClient)

	grpcServer := grpc.NewServer()
	userpb.RegisterUserServiceServer(grpcServer, grpcUserService)
	grpcServer.Serve(lis)
}
