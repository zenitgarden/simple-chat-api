package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/zenitgarden/simple-chat-api/internal/entity"
	"gorm.io/gorm"
)

type ParticipantRepository interface {
	Create(ctx context.Context, participant *entity.Participant, trx *gorm.DB) error
	CreateMany(ctx context.Context, participants []entity.Participant, trx *gorm.DB) error
	FindByUserID(ctx context.Context, userID uuid.UUID, trx *gorm.DB) (*entity.Participant, error)
	FindById(ctx context.Context, conversationID, userId uuid.UUID, trx *gorm.DB) (*entity.Participant, error)
	FindByConversationId(ctx context.Context, conversationId uuid.UUID) ([]*uuid.UUID, error)
	Delete(ctx context.Context, participant *entity.Participant, trx *gorm.DB) error
	BeginTrx(ctx context.Context) *gorm.DB
}

type participantRepository struct {
	db *gorm.DB
}

func NewParticipantRepository(db *gorm.DB) ParticipantRepository {
	return &participantRepository{db}
}

func (r *participantRepository) Create(ctx context.Context, participant *entity.Participant, trx *gorm.DB) error {
	db := r.db
	if trx != nil {
		db = trx
	}

	if err := db.WithContext(ctx).Create(participant).Error; err != nil {
		return nil
	}

	return nil
}

func (r *participantRepository) CreateMany(ctx context.Context, participants []entity.Participant, trx *gorm.DB) error {
	db := r.db
	if trx != nil {
		db = trx
	}

	if err := db.WithContext(ctx).Create(&participants).Error; err != nil {
		return err
	}

	return nil
}

func (r *participantRepository) FindByUserID(ctx context.Context, userID uuid.UUID, trx *gorm.DB) (*entity.Participant, error) {
	var participant entity.Participant

	db := r.db
	if trx != nil {
		db = trx
	}

	if err := db.WithContext(ctx).First(&participant, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}

	return &participant, nil
}

func (r *participantRepository) FindById(ctx context.Context, conversationID, userId uuid.UUID, trx *gorm.DB) (*entity.Participant, error) {
	var participant entity.Participant

	db := r.db
	if trx != nil {
		db = trx
	}

	if err := db.WithContext(ctx).First(&participant, "conversation_id = ? AND user_id = ?", conversationID, userId).Error; err != nil {
		return nil, err
	}

	return &participant, nil
}

func (r *participantRepository) Delete(ctx context.Context, participant *entity.Participant, trx *gorm.DB) error {
	db := r.db
	if trx != nil {
		db = trx
	}

	if err := db.WithContext(ctx).Delete(participant, participant.ID).Error; err != nil {
		return err
	}

	return nil
}

func (r *participantRepository) FindByConversationId(ctx context.Context, conversationId uuid.UUID) ([]*uuid.UUID, error) {
	var participants []struct {
		UserID uuid.UUID
	}

	err := r.db.WithContext(ctx).Model(&entity.Participant{}).Where("conversation_id = ?", conversationId).Find(&participants).Error
	if err != nil {
		return nil, err
	}

	var ids []*uuid.UUID
	for _, p := range participants {
		ids = append(ids, &p.UserID)
	}
	return ids, nil
}

func (r *participantRepository) BeginTrx(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Begin()
}
