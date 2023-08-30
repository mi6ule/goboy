package grpc_main

import (
	"net"

	userpb "gitlab.avakatan.ir/boilerplates/go-boiler/gen/go/proto/user/v1"
	"gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/database/persistence"
	errorhandler "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/error-handler"
	grpc_service "gitlab.avakatan.ir/boilerplates/go-boiler/internal/infrastructure/grpc/service"
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
