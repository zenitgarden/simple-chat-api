package handler

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zenitgarden/simple-chat-api/internal/dto"
	"github.com/zenitgarden/simple-chat-api/internal/service"
	"github.com/zenitgarden/simple-chat-api/pkg/exception"
	"github.com/zenitgarden/simple-chat-api/pkg/utils"
	"github.com/zenitgarden/simple-chat-api/pkg/validator"
)

type ConversationHandler struct {
	conversationService *service.ConversationService
	messageService      *service.MessageService
}

func NewConversationHandler(conversationService *service.ConversationService, messageService *service.MessageService) *ConversationHandler {
	return &ConversationHandler{
		conversationService,
		messageService,
	}
}

func (h *ConversationHandler) CreateConversation(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	id, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.Res{
			StatusCode: fiber.StatusUnauthorized,
			Message:    "Invalid or expired token",
			Data:       nil,
		})
	}

	var req dto.CreateConversationRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		return exception.JSON(c, err)
	}

	if id == req.UserId {
		return c.Status(fiber.StatusBadRequest).JSON(utils.Res{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Cannot create conversation with yourself",
			Data:       nil,
		})
	}

	conversation, err := h.conversationService.CreateConversation(c.Context(), &req, id)

	if err != nil {
		return exception.JSON(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(utils.Res{
		StatusCode: fiber.StatusCreated,
		Message:    "Conversation created successfully",
		Data: dto.ConversationResponse{
			ID:        conversation.ID,
			Title:     conversation.Title,
			IsGroup:   conversation.IsGroup,
			CreatedAt: conversation.CreatedAt,
			UpdatedAt: conversation.UpdatedAt,
		},
	})
}

func (h *ConversationHandler) FindConversationByID(c *fiber.Ctx) error {
	idParam, err := validator.ValidateUUID(c, c.Params("id"))
	if err != nil {
		return exception.JSON(c, err)
	}

	userId := c.Locals("userId").(string)
	userID, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.Res{
			StatusCode: fiber.StatusUnauthorized,
			Message:    "Invalid or expired token",
			Data:       nil,
		})
	}

	conversation, err := h.conversationService.FindConversationByID(c.Context(), idParam, userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Res{
			StatusCode: fiber.StatusNotFound,
			Message:    "Conversation not found",
			Data:       nil,
		})
	}

	var participants []dto.ParticipantResponse
	for _, participant := range conversation.Participants {
		if !conversation.IsGroup && participant.User.ID != userID {
			conversation.Title = participant.User.Name
		}

		participants = append(participants, dto.ParticipantResponse{
			ID:       participant.User.ID,
			Name:     participant.User.Name,
			JoinedAt: participant.JoinedAt,
		})
	}
	return c.Status(fiber.StatusOK).JSON(utils.Res{
		StatusCode: fiber.StatusOK,
		Message:    "Conversation found",
		Data: dto.ConversationResponse{
			ID:           conversation.ID,
			Title:        conversation.Title,
			IsGroup:      conversation.IsGroup,
			CreatedAt:    conversation.CreatedAt,
			UpdatedAt:    conversation.UpdatedAt,
			Participants: &participants,
		},
	})
}
func (h *ConversationHandler) UpdateConversation(c *fiber.Ctx) error {
	idParam, err := validator.ValidateUUID(c, c.Params("id"))
	if err != nil {
		return exception.JSON(c, err)
	}

	userId := c.Locals("userId").(string)
	userID, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.Res{
			StatusCode: fiber.StatusUnauthorized,
			Message:    "Invalid or expired token",
			Data:       nil,
		})
	}

	var req dto.UpdateConversationRequest
	if err := validator.ValidateRequest(c, &req); err != nil {
		return exception.JSON(c, err)
	}

	conversation, err := h.conversationService.UpdateConversation(c.Context(), idParam, userID, &req)
	if err != nil {
		return exception.JSON(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(utils.Res{
		StatusCode: fiber.StatusOK,
		Message:    "Conversation updated successfully",
		Data: dto.ConversationResponse{
			ID:        conversation.ID,
			Title:     conversation.Title,
			IsGroup:   conversation.IsGroup,
			CreatedAt: conversation.CreatedAt,
			UpdatedAt: conversation.UpdatedAt,
		},
	})
}

func (h *ConversationHandler) FindAllConversations(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	id, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.Res{
			StatusCode: fiber.StatusUnauthorized,
			Message:    "Invalid or expired token",
			Data:       nil,
		})
	}

	allowedColumn := []string{"title", "is_group", "created_at"}

	pagination := utils.GetPagination(c)
	filter := dto.ConversationFilter{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
		Title:  c.Query("name"),
		Sort:   utils.ParseSort(c.Query("sort"), allowedColumn, "created_at desc"),
	}

	conversations, total, err := h.conversationService.FindAllConversations(c.Context(), filter, id)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(utils.Res{
			StatusCode: fiber.StatusServiceUnavailable,
			Message:    "Service unavailable",
			Data:       nil,
		})
	}

	var result []dto.AllConversationResponse = []dto.AllConversationResponse{}
	for _, conversation := range conversations {
		var sentAt *time.Time
		if !conversation.MessageSentAt.IsZero() {
			sentAt = &conversation.MessageSentAt
		} else {
			sentAt = nil
		}
		result = append(result, dto.AllConversationResponse{
			ID: conversation.ConversationID,
			Title: func() string {
				if conversation.IsGroup {
					return conversation.Title
				} else {
					if conversation.ParticipantName != "" {
						return conversation.ParticipantName
					}
					return conversation.Title
				}
			}(),
			IsGroup:     conversation.IsGroup,
			LastMessage: conversation.LatestMessage,
			SentAt:      sentAt,
			CreatedBy:   conversation.CreatedBy,
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

func (h *ConversationHandler) GetMessagesByConversationID(c *fiber.Ctx) error {
	idParam, err := validator.ValidateUUID(c, c.Params("id"))
	if err != nil {
		return exception.JSON(c, err)
	}

	userId := c.Locals("userId").(string)
	userID, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.Res{
			StatusCode: fiber.StatusUnauthorized,
			Message:    "Invalid or expired token",
			Data:       nil,
		})
	}

	conversation, err := h.conversationService.FindByID(c.Context(), idParam, userID)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(utils.Res{
			StatusCode: fiber.StatusServiceUnavailable,
			Message:    "Service unavailable",
			Data:       nil,
		})
	}

	if conversation == nil {
		return c.Status(fiber.StatusNotFound).JSON(utils.Res{
			StatusCode: fiber.StatusNotFound,
			Message:    "Conversation not found",
			Data:       nil,
		})
	}

	pagination := utils.GetPagination(c)
	filter := dto.MessageFilter{
		ConversationID: conversation.ID,
		Limit:          pagination.Limit,
		Offset:         pagination.Offset,
	}

	messages, total, err := h.messageService.FindAllByConversationID(c.Context(), filter)
	if err != nil {
		return exception.JSON(c, err)
	}

	var result []dto.MessageResponse = []dto.MessageResponse{}
	for _, message := range messages {
		result = append(result, dto.MessageResponse{
			ID:             message.ID,
			ConversationID: message.ConversationID,
			SenderID:       message.SenderID,
			SenderName:     message.User.Name,
			Content:        message.Content,
			SentAt:         message.SentAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.Res{
		StatusCode: fiber.StatusOK,
		Message:    "Messages found successfully",
		Data:       result,
		TotalData:  total,
		Pagination: &utils.Pagination{
			Page:  pagination.Page,
			Limit: pagination.Limit,
		},
	})
}

func (h *ConversationHandler) GetLatestConversation(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)
	userID, err := uuid.Parse(userId)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(utils.Res{
			StatusCode: fiber.StatusUnauthorized,
			Message:    "Invalid or expired token",
			Data:       nil,
		})
	}

	conversation, err := h.conversationService.GetLatestConversation(c.Context(), userID)
	if err != nil {

		if err.Error() == "record not found" {
			return c.Status(fiber.StatusOK).JSON(utils.Res{
				StatusCode: fiber.StatusOK,
				Message:    "No conversation found",
				Data:       conversation,
			})
		}

		return c.Status(fiber.StatusServiceUnavailable).JSON(utils.Res{
			StatusCode: fiber.StatusServiceUnavailable,
			Message:    "Service unavailable",
			Data:       nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(utils.Res{
		StatusCode: fiber.StatusOK,
		Message:    "Success",
		Data:       conversation,
	})
}
