package service

import (
	"context"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/zenitgarden/simple-chat-api/internal/dto"
	"github.com/zenitgarden/simple-chat-api/internal/entity"
	"github.com/zenitgarden/simple-chat-api/internal/repository"
	"github.com/zenitgarden/simple-chat-api/pkg/exception"
)

type ConversationService struct {
	conversationRepo repository.ConversationRepository
	participantRepo  repository.ParticipantRepository
}

func NewConversationService(conversationRepo repository.ConversationRepository, participantRepo repository.ParticipantRepository) *ConversationService {
	return &ConversationService{
		conversationRepo: conversationRepo,
		participantRepo:  participantRepo,
	}
}
func (s *ConversationService) CreateConversation(ctx context.Context, req *dto.CreateConversationRequest, userID uuid.UUID) (*entity.Conversation, error) {

	conversation := &entity.Conversation{
		ID:        uuid.New(),
		Title:     req.Title,
		IsGroup:   req.IsGroup,
		CreatedBy: userID,
	}

	tx := s.conversationRepo.BeginTrx(ctx)
	if tx.Error != nil {
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to create conversation",
			Data:       nil,
		}
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	participants := []entity.Participant{
		{
			ID:             uuid.New(),
			ConversationID: conversation.ID,
			UserID:         userID,
		},
	}

	if !req.IsGroup && req.UserId != uuid.Nil {
		exists, err := s.conversationRepo.CheckConversation(ctx, req.UserId, userID, tx)

		fmt.Println("exists", exists)

		if exists == 1 {
			tx.Rollback()
			return nil, &exception.HTTPError{
				StatusCode: fiber.StatusConflict,
				Message:    "Conversation already exists",
				Data:       nil,
			}
		}
		if err != nil {
			tx.Rollback()
			return nil, &exception.HTTPError{
				StatusCode: fiber.StatusServiceUnavailable,
				Message:    "Failed to create conversation",
				Data:       nil,
			}
		}

		participants = append(participants, entity.Participant{
			ID:             uuid.New(),
			ConversationID: conversation.ID,
			UserID:         req.UserId,
		})
	}

	err := s.conversationRepo.Create(ctx, conversation, tx)
	if err != nil {
		tx.Rollback()
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusServiceUnavailable,
			Message:    "Failed to create conversation",
			Data:       nil,
		}
	}

	err = s.participantRepo.CreateMany(ctx, participants, tx)

	if err != nil {
		fmt.Println("exists 3", err)
		tx.Rollback()
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusServiceUnavailable,
			Message:    "Failed to create conversation",
			Data:       nil,
		}
	}

	return conversation, tx.Commit().Error
}
func (s *ConversationService) FindConversationByID(ctx context.Context, id, userId uuid.UUID) (*entity.Conversation, error) {
	participant, err := s.participantRepo.FindById(ctx, id, userId, nil)

	if err != nil {
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusNotFound,
			Message:    "Conversation not found",
			Data:       nil,
		}
	}

	if participant == nil {
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusNotFound,
			Message:    "Conversation not found",
			Data:       nil,
		}
	}

	return s.conversationRepo.FindByIDPreload(ctx, id, nil)
}

func (s *ConversationService) FindByID(ctx context.Context, id, userId uuid.UUID) (*entity.Conversation, error) {
	return s.conversationRepo.CheckParticipant(ctx, id, userId, nil)
}

func (s *ConversationService) UpdateConversation(ctx context.Context, id, userId uuid.UUID, req *dto.UpdateConversationRequest) (*entity.Conversation, error) {
	tx := s.conversationRepo.BeginTrx(ctx)
	if tx.Error != nil {
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to update conversation",
			Data:       nil,
		}
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	conversation, err := s.conversationRepo.CheckParticipant(ctx, id, userId, tx)

	if err != nil {
		tx.Rollback()
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusNotFound,
			Message:    "Conversation not found",
			Data:       nil,
		}
	}

	if req.Title != "" {
		conversation.Title = req.Title
	}

	if req.IsGroup != nil {
		conversation.IsGroup = *req.IsGroup
	}

	err = s.conversationRepo.Update(ctx, conversation, tx)
	if err != nil {
		tx.Rollback()
		return nil, &exception.HTTPError{
			StatusCode: fiber.StatusServiceUnavailable,
			Message:    "Failed to update conversation",
			Data:       nil,
		}
	}

	return conversation, tx.Commit().Error
}

func (s *ConversationService) FindAllConversations(ctx context.Context, filter dto.ConversationFilter, userId uuid.UUID) ([]*entity.ConversationSummary, int64, error) {
	return s.conversationRepo.FindAll(ctx, filter, userId)
}

func (s *ConversationService) GetLatestConversation(ctx context.Context, userId uuid.UUID) (*dto.ConversationDetail, error) {
	return s.conversationRepo.GetLatestConversation(ctx, userId)
}

func (s *ConversationService) FindGroupConversation(ctx context.Context, filter dto.ConversationFilter, userId uuid.UUID) ([]dto.GroupConversation, int64, error) {
	return s.conversationRepo.FindGroupConversation(ctx, filter, userId)
}
