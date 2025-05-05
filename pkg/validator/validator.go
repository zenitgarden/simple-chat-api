package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zenitgarden/simple-chat-api/pkg/exception"
)

var validate = validator.New()

func ValidateRequest(c *fiber.Ctx, req any) error {
	if err := c.BodyParser(req); err != nil {
		if strings.Contains(err.Error(), "invalid UUID") {
			return &exception.HTTPError{
				StatusCode: fiber.StatusBadRequest,
				Message:    "Invalid UUID",
				Data:       nil,
			}
		}
		return &exception.HTTPError{
			StatusCode: fiber.StatusUnprocessableEntity,
			Message:    "Unprocessable entity",
			Data:       nil,
		}
	}

	if err := validate.Struct(req); err != nil {
		var messages string
		for _, e := range err.(validator.ValidationErrors) {
			switch e.Tag() {
			case "required":
				messages = e.Field() + " is required"
			case "email":
				messages = e.Field() + " must be a valid email"
			case "min":
				messages = e.Field() + " must be at least " + e.Param() + " characters"
			case "max":
				messages = e.Field() + " must be at most " + e.Param() + " characters"
			default:
				messages = e.Field() + " is invalid"
			}
		}
		return &exception.HTTPError{
			StatusCode: fiber.StatusBadRequest,
			Message:    messages,
			Data:       nil,
		}
	}
	return nil
}
func ValidateUUID(c *fiber.Ctx, id string) (uuid.UUID, error) {
	idParam, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return uuid.UUID{}, &exception.HTTPError{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid UUID",
			Data:       nil,
		}
	}
	return idParam, nil
}
