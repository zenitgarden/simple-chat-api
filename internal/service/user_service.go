package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zenitgarden/simple-chat-api/internal/dto"
	"github.com/zenitgarden/simple-chat-api/internal/entity"
	"github.com/zenitgarden/simple-chat-api/internal/repository"
	"github.com/zenitgarden/simple-chat-api/pkg/exception"
	"github.com/zenitgarden/simple-chat-api/pkg/security"
)

type UserService struct {
	UserRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

func (s *UserService) RegisterUser(ctx context.Context, req *dto.RegisterRequest) (*entity.User, error) {
	hashedPassword, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusServiceUnavailable,
			Message:    "Service unavailable",
			Data:       nil,
		}
	}

	user := &entity.User{
		ID:       uuid.New(),
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	return user, s.UserRepo.Create(ctx, user)
}

func (s *UserService) FindUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	return s.UserRepo.FindByID(ctx, id, nil)
}

func (s *UserService) FindAll(ctx context.Context, filter dto.UserFilter) ([]*entity.User, int64, error) { return s.UserRepo.FindAll(ctx, filter) }

func (s *UserService) UpdateUser(ctx context.Context, id uuid.UUID, req *dto.UpdateUserRequest) (*entity.User, error) {

	// Begin transaction
	tx := s.UserRepo.BeginTrx(ctx)
	if tx.Error != nil {
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to update user",
			Data:       nil,
		}
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user, err := s.UserRepo.FindByID(ctx, id, tx)
	if err != nil {
		tx.Rollback()
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusNotFound,
			Message:    "User not found",
			Data:       nil,
		}
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	err = s.UserRepo.Update(ctx, user, tx)
	if err != nil {
		tx.Rollback()
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to update user",
			Data:       nil,
		}
	}

	return user, tx.Commit().Error
}
