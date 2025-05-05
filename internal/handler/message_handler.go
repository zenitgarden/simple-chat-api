package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/zenitgarden/simple-chat-api/internal/dto"
	"github.com/zenitgarden/simple-chat-api/internal/service"
	"github.com/zenitgarden/simple-chat-api/pkg/exception"
	"github.com/zenitgarden/simple-chat-api/pkg/utils"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) FindMessagesByConversationID(c *fiber.Ctx) error {
	pagination := utils.GetPagination(c)
	filter := dto.MessageFilter{
		Limit:  pagination.Limit,
		Offset: pagination.Offset,
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
