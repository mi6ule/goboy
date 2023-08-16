package grpc_service

import (
	"context"
	"fmt"

	userpb "gitlab.avakatan.ir/boilerplates/go-boiler/gen/go/proto/user/v1"
)

type GrpcUserService struct {
}

// this type implements the type bellwo from user_service.pb.go
// type UserServiceServer interface {
// 	GetUser(context.Context, *GetUserRequest) (*GetUserResponse, error)
// 	mustEmbedUnimplementedUserServiceServer()
// }

func (u *GrpcUserService) GetUser(_ context.Context, req *userpb.GetUserRequest) (*userpb.GetUserResponse, error) {
	fmt.Printf(req.Username)
	return &userpb.GetUserResponse{User: &userpb.User{Id: "someId"}}, nil
}
