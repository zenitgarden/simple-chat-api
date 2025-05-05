package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zenitgarden/simple-chat-api/internal/dto"
	"github.com/zenitgarden/simple-chat-api/internal/service"
	"github.com/zenitgarden/simple-chat-api/pkg/exception"
	"github.com/zenitgarden/simple-chat-api/pkg/utils"
	"github.com/zenitgarden/simple-chat-api/pkg/validator"
)

type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := validator.ValidateRequest(c, &req); err != nil {
		return exception.JSON(c, err)
	}

	user, token, err := h.authService.Login(c.Context(), req.Email, req.Password)
	if err != nil {
		return exception.JSON(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.Res{
		StatusCode: fiber.StatusOK,
		Message:    "Login successful",
		Data: dto.LoginResponse{
			UserId:     user.ID,
			UserName:   user.Name,
			AccesToken: *token,
		},
	})
}

func (h *AuthHandler) RegisterUser(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := validator.ValidateRequest(c, &req); err != nil {
		return exception.JSON(c, err)
	}

	// Register the user
	user, err := h.userService.RegisterUser(c.Context(), &req)
	if err != nil {
		if err.Error() == "email" {
			return c.Status(fiber.StatusConflict).JSON(utils.Res{
				StatusCode: fiber.StatusConflict,
				Message:    "Email already exists",
				Data:       nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Res{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to register user",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Res{
		StatusCode: fiber.StatusCreated,
		Message:    "User registered successfully",
		Data: dto.UserResponseDetail{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	id, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.Res{
			StatusCode: fiber.StatusUnauthorized,
			Message:    "Invalid or expired token",
			Data:       nil,
		})
	}

	user, err := h.userService.FindUserByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.Res{
			StatusCode: fiber.StatusUnauthorized,
			Message:    "Invalid or expired token",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.Res{
		StatusCode: fiber.StatusOK,
		Message:    "Success",
		Data: dto.UserResponseDetail{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}
