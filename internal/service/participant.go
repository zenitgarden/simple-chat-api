package service

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zenitgarden/simple-chat-api/internal/entity"
	"github.com/zenitgarden/simple-chat-api/internal/repository"
	"github.com/zenitgarden/simple-chat-api/pkg/exception"
)

type ParticipantService struct {
	participantRepository  repository.ParticipantRepository
	conversationRepository repository.ConversationRepository
}

func NewParticipantService(participantRepository repository.ParticipantRepository, conversationRepository repository.ConversationRepository) *ParticipantService {
	return &ParticipantService{
		participantRepository:  participantRepository,
		conversationRepository: conversationRepository,
	}
}

func (s *ParticipantService) CreateParticipant(ctx context.Context, conversationID, userID uuid.UUID) (*entity.Participant, error) {
	tx := s.participantRepository.BeginTrx(ctx)
	if tx.Error != nil {
		tx.Rollback()
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to create participant",
			Data:       nil,
		}
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	participant := &entity.Participant{
		ConversationID: conversationID,
		UserID:         userID,
	}

	conversation, err := s.conversationRepository.FindByID(ctx, conversationID, tx)
	if conversation == nil {
		tx.Rollback()
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusNotFound,
			Message:    "Conversation not found",
			Data:       nil,
		}
	}

	if err != nil {
		tx.Rollback()
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusServiceUnavailable,
			Message:    "Failed to create participant",
			Data:       nil,
		}
	}

	err = s.participantRepository.Create(ctx, participant, tx)
	if err != nil {
		tx.Rollback()
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to create participant",
			Data:       nil,
		}
	}
	return participant, tx.Commit().Error
}

func (s *ParticipantService) DeleteParticipant(ctx context.Context, id uuid.UUID) error {
	participant := &entity.Participant{
		ID: id,
	}

	tx := s.participantRepository.BeginTrx(ctx)
	if tx.Error != nil {
		tx.Rollback()
		return &exception.HTTPError{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to delete participant",
			Data:       nil,
		}
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	err := s.participantRepository.Delete(ctx, participant, tx)
	if err != nil {
		tx.Rollback()
		return &exception.HTTPError{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to delete participant",
			Data:       nil,
		}
	}
	return tx.Commit().Error
}

func (s *ParticipantService) FindByConversationId(ctx context.Context, id uuid.UUID) ([]*uuid.UUID, error) {
	return s.participantRepository.FindByConversationId(ctx, id)
}
