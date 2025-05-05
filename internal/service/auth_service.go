package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/zenitgarden/simple-chat-api/internal/entity"
	"github.com/zenitgarden/simple-chat-api/internal/repository"
	"github.com/zenitgarden/simple-chat-api/pkg/exception"
	"github.com/zenitgarden/simple-chat-api/pkg/security"
	"github.com/zenitgarden/simple-chat-api/pkg/utils"
)

type AuthService struct {
	UserRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*entity.User, *string, error) {
	user, err := s.UserRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, nil, &exception.HTTPError{
			StatusCode: fiber.StatusUnauthorized,
			Message:    "Email or password is incorrect",
			Data:       nil,
		}
	}

	if !security.CheckPasswordHash(password, user.Password) {
		return nil, nil, &exception.HTTPError{
			StatusCode: fiber.StatusUnauthorized,
			Message:    "Email or password is incorrect",
			Data:       nil,
		}
	}

	token, err := utils.GenerateJWT(user)
	if err != nil {
		return nil, nil, &exception.HTTPError{
			StatusCode: fiber.StatusServiceUnavailable,
			Message:    "Failed to login",
			Data:       nil,
		}
	}

	return user, token, nil
}
