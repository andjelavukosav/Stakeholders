package handlers

import (
	"context"
	"database-example/model"
	pb "database-example/proto/stakeholders"
	"database-example/service"
	"database-example/util"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserHandler struct {
	pb.UnimplementedStakeholdersServiceServer
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (h *UserHandler) RegisterUser(ctx context.Context, req *pb.RegistrationRequest) (*pb.AuthenticationResponse, error) {
	user := &model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		Role:     req.Role,
	}

	if err := h.UserService.CreateUser(user); err != nil {
		return nil, status.Errorf(codes.Internal, "could not create user: %v", err)
	}

	token, err := util.GenerateToken(user.ID.String(), user.Username, user.Role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &pb.AuthenticationResponse{Token: token}, nil
}

func (h *UserHandler) LoginUser(ctx context.Context, req *pb.LoginRequest) (*pb.AuthenticationResponse, error) {
	user, err := h.UserService.Authenticate(req.Username, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	token, err := util.GenerateToken(user.ID.String(), user.Username, user.Role)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token")
	}

	return &pb.AuthenticationResponse{Token: token}, nil
}
