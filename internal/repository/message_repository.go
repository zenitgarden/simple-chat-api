package repository

import (
	"context"

	"github.com/zenitgarden/simple-chat-api/internal/dto"
	"github.com/zenitgarden/simple-chat-api/internal/entity"
	"gorm.io/gorm"
)

type MessageRepository interface {
	Create(ctx context.Context, message *entity.Message) error
	FindAllByConversationID(ctx context.Context, filter dto.MessageFilter) ([]*entity.Message, int64, error)
}

type messageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db}
}

func (r *messageRepository) Create(ctx context.Context, message *entity.Message) error {
	// Start the transaction
	tx := r.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Insert message
	if err := tx.Create(message).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}

func (r *messageRepository) FindAllByConversationID(ctx context.Context, filter dto.MessageFilter) ([]*entity.Message, int64, error) {
	var messages []*entity.Message
	var count int64

	query := r.db.WithContext(ctx).
		Model(&entity.Message{}).
		Where("conversation_id = ?", filter.ConversationID)

	// Get total count
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Apply ordering, pagination, and retrieve messages
	if err := query.
		Preload("User").
		Order("sent_at desc").
		Limit(filter.Limit).
		Offset(filter.Offset).
		Find(&messages).Error; err != nil {
		return nil, 0, err
	}
	return messages, count, nil
}
