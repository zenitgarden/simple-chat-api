package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zenitgarden/simple-chat-api/internal/dto"
	"github.com/zenitgarden/simple-chat-api/internal/service"
	"github.com/zenitgarden/simple-chat-api/pkg/exception"
	"github.com/zenitgarden/simple-chat-api/pkg/utils"
	"github.com/zenitgarden/simple-chat-api/pkg/validator"
)

type ParticipantHandler struct {
	participantService *service.ParticipantService
}

func NewParticpantHandler(participantService *service.ParticipantService) *ParticipantHandler {
	return &ParticipantHandler{participantService: participantService}
}

func (h *ParticipantHandler) CreateParticipant(c *fiber.Ctx) error {
	
	var req dto.CreateParticipantRequest

	if err := validator.ValidateRequest(c, &req); err != nil {
		return exception.JSON(c, err)
	}

	participant, err := h.participantService.CreateParticipant(c.Context(), req.ConversationID, req.UserID)

	if err != nil {
		return exception.JSON(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Res{
		StatusCode: fiber.StatusCreated,
		Message:    "Participant created successfully",
		Data: dto.ParticipantRawResponse{
			ConversationID: participant.ConversationID,
			UserID:         participant.UserID,
			JoinedAt:       participant.JoinedAt,
		},
	})
}

func (h *ParticipantHandler) DeleteParticipant(c *fiber.Ctx) error {
	idParam, err := validator.ValidateUUID(c, c.Params("id"))

	if err != nil {
		return exception.JSON(c, err)
	}

	err = h.participantService.DeleteParticipant(c.Context(), idParam)

	if err != nil {
		return exception.JSON(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Res{
		StatusCode: fiber.StatusCreated,
		Message:    "Participant deleted successfully",
		Data:       nil,
	})
}
