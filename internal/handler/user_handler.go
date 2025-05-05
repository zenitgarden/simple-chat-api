// internal/handler/user_handler.go
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

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterUser(c *fiber.Ctx) error {
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

func (h *UserHandler) FindUserByID(c *fiber.Ctx) error {
	idParam, err := validator.ValidateUUID(c, c.Params("id"))
	if err != nil {
		return exception.JSON(c, err)
	}

	user, err := h.userService.FindUserByID(c.Context(), idParam)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Res{
			StatusCode: fiber.StatusNotFound,
			Message:    "User not found",
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

func (h *UserHandler) FindUsers(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	id, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.Res{
			StatusCode: fiber.StatusUnauthorized,
			Message:    "Invalid or expired token",
			Data:       nil,
		})
	}

	pagination := utils.GetPagination(c)
	filter := dto.UserFilter{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
		Name:   c.Query("name"),
		UserID: id,
	}

	users, total, err := h.userService.FindAll(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(utils.Res{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to find users",
			Data:       nil,
		})
	}

	var result []dto.UserResponse = []dto.UserResponse{}
	for _, user := range users {
		result = append(result, dto.UserResponse{
			ID:   user.ID,
			Name: user.Name,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.Res{
		StatusCode: fiber.StatusOK,
		Message:    "Success",
		Data:       result,
		TotalData:  total,
		Pagination: &utils.Pagination{
			Page:  pagination.Page,
			Limit: pagination.Limit,
		},
	})
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idParam, err := validator.ValidateUUID(c, c.Params("id"))
	if err != nil {
		return exception.JSON(c, err)
	}

	var req dto.UpdateUserRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		return exception.JSON(c, err)
	}

	user, err := h.userService.UpdateUser(c.Context(), idParam, &req)
	if err != nil {
		return exception.JSON(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.Res{
		StatusCode: fiber.StatusOK,
		Message:    "User updated successfully",
		Data: dto.UserResponseDetail{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}
