package service

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zenitgarden/simple-chat-api/internal/dto"
	"github.com/zenitgarden/simple-chat-api/internal/entity"
	"github.com/zenitgarden/simple-chat-api/internal/repository"
	"github.com/zenitgarden/simple-chat-api/pkg/exception"
)

type MessageService struct {
	messageRepo repository.MessageRepository
}

func NewMessageService(messageRepo repository.MessageRepository) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}

func (s *MessageService) CreateMessage(ctx context.Context, req *dto.CreateMessageRequest) (*entity.Message, error) {
	message := &entity.Message{
		ID:             uuid.New(),
		ConversationID: req.ConversationID,
		SenderID:       req.SenderID,
		Content:        req.Content,
		SentAt:         time.Now(),
	}

	err := s.messageRepo.Create(ctx, message)
	if err != nil {
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to create message",
			Data:       nil,
		}
	}

	return message, nil
}

func (s *MessageService) FindAllByConversationID(ctx context.Context, filter dto.MessageFilter) ([]*entity.Message, int64, error) {
	return s.messageRepo.FindAllByConversationID(ctx, filter)
}
