package handlers

import (
	"context"
	pb "database-example/proto/stakeholders"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Samo implementiraš metodu na istom handleru.
func (h *UserHandler) GetAllUsers(ctx context.Context, req *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	users, err := h.UserService.GetAllUsers()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	var pbUsers []*pb.User
	for _, u := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:        u.ID.String(),
			Username:  u.Username,
			Email:     u.Email,
			Role:      u.Role,
			IsBlocked: u.IsBlocked,
		})
	}

	return &pb.GetAllUsersResponse{
		Users: pbUsers,
	}, nil
}

func (h *UserHandler) BlockUser(ctx context.Context, req *pb.BlockUserRequest) (*pb.BlockUserResponse, error) {
	// pozivamo servisnu logiku
	err := h.UserService.BlockUser(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to block user: %v", err)
	}

	return &pb.BlockUserResponse{
		Message: "User blocked successfully",
	}, nil
}

func (h *UserHandler) UnblockUser(ctx context.Context, req *pb.UnblockUserRequest) (*pb.UnblockUserResponse, error) {
	// pozivamo servisnu logiku
	err := h.UserService.UnblockUser(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to unblock user: %v", err)
	}

	return &pb.UnblockUserResponse{
		Message: "User unblocked successfully",
	}, nil
}
