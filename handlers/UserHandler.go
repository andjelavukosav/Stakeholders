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

func (h *UserHandler) GetUserProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.UserProfileResponse, error) {
	user, err := h.UserService.GetUserProfile(req.Username)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &pb.UserProfileResponse{
		Id:           user.ID.String(),
		Username:     user.Username,
		Email:        user.Email,
		Role:         user.Role,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		ProfileImage: user.ProfileImage,
		Biography:    user.Biography,
		Motto:        user.Motto,
	}, nil
}

func (h *UserHandler) UpdateUserProfile(ctx context.Context, req *pb.UpdateProfileRequest) (*pb.UpdateProfileResponse, error) {
	updatedUser := &model.User{
		Username:     req.Username,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		ProfileImage: req.ProfileImage,
		Biography:    req.Biography,
		Motto:        req.Motto,
	}

	err := h.UserService.UpdateUserProfile(req.Username, updatedUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update profile: %v", err)
	}

	return &pb.UpdateProfileResponse{
		Message: "Profile updated successfully",
		Success: true,
	}, nil
}

func (h *UserHandler) UploadProfileImage(ctx context.Context, req *pb.UploadProfileImageRequest) (*pb.UploadProfileImageResponse, error) {
    // validacija
    if req.Username == "" || req.ImagePath == "" {
        return nil, status.Errorf(codes.InvalidArgument, "username and image path are required")
    }

    // poziv servisa da upiše image url u bazu
    err := h.UserService.UploadProfileImage(req.Username, req.ImagePath)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to update profile image: %v", err)
    }

    return &pb.UploadProfileImageResponse{
        Success: true,
        Message: "Profile image uploaded successfully",
        ImageUrl: req.ImagePath,
    }, nil
}



