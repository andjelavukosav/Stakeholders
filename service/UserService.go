package service

import (
	"context"
	"database-example/model"
	"database-example/repo"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	UserRepo *repo.UserRepository
}

// Registracija
func (s *UserService) CreateUser(user *model.User) error {
	// Hash lozinke pre upisa u bazu
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)

	user.BeforeCreate()
	ctx := context.Background()
	return s.UserRepo.CreateUser(ctx, user)
}

// Login
func (s *UserService) Authenticate(username, password string) (*model.User, error) {
	ctx := context.Background()
	user, err := s.UserRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Uporedi password sa hash-om iz baze
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if user.IsBlocked {
		return nil, errors.New("user is blocked")
	}

	return user, nil
}

// Dohvati user-a po username
func (s *UserService) GetByUsername(username string) (*model.User, error) {
	ctx := context.Background()
	return s.UserRepo.FindByUsername(ctx, username)
}

// Dohvati sve korisnike
func (s *UserService) GetAllUsers() ([]*model.User, error) {
	ctx := context.Background()
	return s.UserRepo.GetAllUsers(ctx)
}

func (s *UserService) BlockUser(userID string) error {
	ctx := context.Background()
	return s.UserRepo.SetUserBlocked(ctx, userID, true)
}

func (s *UserService) UnblockUser(userID string) error {
	ctx := context.Background()
	return s.UserRepo.SetUserBlocked(ctx, userID, false)
}
