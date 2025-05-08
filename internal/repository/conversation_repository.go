package repository

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/zenitgarden/simple-chat-api/internal/dto"
	"github.com/zenitgarden/simple-chat-api/internal/entity"
	"gorm.io/gorm"
)

type ConversationRepository interface {
	Create(ctx context.Context, conversation *entity.Conversation, trx *gorm.DB) error
	FindByIDPreload(ctx context.Context, id uuid.UUID, trx *gorm.DB) (*entity.Conversation, error)
	CheckParticipant(ctx context.Context, id, userId uuid.UUID, trx *gorm.DB) (*entity.Conversation, error)
	CheckConversation(ctx context.Context, user1, user2 uuid.UUID, trx *gorm.DB) (int8, error)
	FindByID(ctx context.Context, id uuid.UUID, trx *gorm.DB) (*entity.Conversation, error)
	Update(ctx context.Context, conversation *entity.Conversation, trx *gorm.DB) error
	FindAll(ctx context.Context, filter dto.ConversationFilter, userId uuid.UUID) ([]*entity.ConversationSummary, int64, error)
	GetLatestConversation(ctx context.Context, userId uuid.UUID) (*dto.ConversationDetail, error)
	FindGroupConversation(ctx context.Context, filter dto.ConversationFilter, userId uuid.UUID) ([]dto.GroupConversation, int64, error)
	BeginTrx(ctx context.Context) *gorm.DB
}

type conversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) ConversationRepository {
	return &conversationRepository{db}
}
func (r *conversationRepository) Create(ctx context.Context, conversation *entity.Conversation, trx *gorm.DB) error {
	db := r.db
	if trx != nil {
		db = trx
	}

	if err := db.WithContext(ctx).Create(conversation).Error; err != nil {
		return nil
	}

	return nil
}
func (r *conversationRepository) FindByIDPreload(ctx context.Context, id uuid.UUID, trx *gorm.DB) (*entity.Conversation, error) {
	var conversation entity.Conversation

	db := r.db
	if trx != nil {
		db = trx
	}

	if err := db.WithContext(ctx).
		Preload("Participants").
		Preload("Participants.User").
		First(&conversation, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (r *conversationRepository) CheckParticipant(ctx context.Context, id, userId uuid.UUID, trx *gorm.DB) (*entity.Conversation, error) {
	var conversation entity.Conversation

	db := r.db
	if trx != nil {
		db = trx
	}

	if err := db.WithContext(ctx).
		Joins("JOIN participants ON participants.conversation_id = conversations.id").
		Where("conversations.id = ? AND participants.user_id = ?", id, userId).
		First(&conversation).Error; err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (r *conversationRepository) FindByID(ctx context.Context, id uuid.UUID, trx *gorm.DB) (*entity.Conversation, error) {
	var conversation entity.Conversation

	db := r.db
	if trx != nil {
		db = trx
	}

	if err := db.WithContext(ctx).
		First(&conversation, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (r *conversationRepository) FindAll(ctx context.Context, filter dto.ConversationFilter, userId uuid.UUID) ([]*entity.ConversationSummary, int64, error) {
	var results []*entity.ConversationSummary
	var total int64

	err := r.db.Raw(`
	SELECT COUNT(DISTINCT c.id) AS total
		FROM conversations c
		LEFT JOIN participants p ON p.conversation_id = c.id AND p.user_id = ?
	WHERE p.user_id IS NOT NULL OR c.created_by = ?`, userId, userId).Scan(&total).Error

	if err != nil {
		return nil, 0, err
	}

	err = r.db.Raw(`
		SELECT 
			c.id AS conversation_id,
			c.title AS title,
			c.is_group AS is_group,
			m.content AS latest_message,
			m.sent_at AS message_sent_at,
			c.created_by
		FROM conversations c
		LEFT JOIN participants p ON p.conversation_id = c.id AND p.user_id = ?
		-- LEFT JOIN LATERAL to include conversations with no messages
		LEFT JOIN LATERAL (
			SELECT content, sent_at
			FROM messages
			WHERE conversation_id = c.id
			ORDER BY sent_at DESC
			LIMIT 1
		) m ON true
		-- Only include conversations where user is a participant or the creator
		WHERE p.user_id IS NOT NULL OR c.created_by = ?

		ORDER BY m.sent_at DESC NULLS LAST
		LIMIT ? OFFSET ?
	`, userId, userId, filter.Limit, filter.Offset).Scan(&results).Error

	if err != nil {
		return nil, 0, err
	}

	var participants []struct {
		ConversationID  uuid.UUID
		ParticipantName string
	}

	var conversationIds []uuid.UUID
	for _, result := range results {
		if !result.IsGroup {
			conversationIds = append(conversationIds, result.ConversationID)
		}
	}

	err = r.db.Raw(`
		SELECT DISTINCT ON (c.id)
			c.id AS conversation_id,
			u.name AS participant_name
		FROM conversations c
		JOIN participants p ON p.conversation_id = c.id
		JOIN users u ON u.id = p.user_id
		WHERE c.id IN ? 
		AND u.id <> ?
  		`, conversationIds, userId).Scan(&participants).Error

	if err != nil {
		return nil, 0, err
	}

	for _, result := range results {
		for _, participant := range participants {
			if result.ConversationID == participant.ConversationID {
				result.ParticipantName = participant.ParticipantName
				break
			}
		}
	}

	return results, total, nil
}

func (r *conversationRepository) Update(ctx context.Context, conversation *entity.Conversation, trx *gorm.DB) error {
	db := r.db
	if trx != nil {
		db = trx
	}

	if err := db.WithContext(ctx).Save(conversation).Error; err != nil {
		return nil
	}

	return nil
}

func (r *conversationRepository) GetLatestConversation(ctx context.Context, userId uuid.UUID) (*dto.ConversationDetail, error) {
	type ConversationIDResult struct {
		ConversationID uuid.UUID
	}

	var result ConversationIDResult

	// Step 1: Get latest conversation ID
	err := r.db.Raw(`
		WITH user_conversations AS (
			SELECT c.id
			FROM conversations c
			JOIN participants p ON c.id = p.conversation_id
			WHERE p.user_id = ?
		),
		latest_conversation AS (
			SELECT m.conversation_id
			FROM messages m
			WHERE m.conversation_id IN (SELECT id FROM user_conversations)
			ORDER BY m.sent_at DESC
			LIMIT 1
		)
		SELECT conversation_id FROM latest_conversation
	`, userId).Scan(&result).Error

	conversationID := result.ConversationID
	if err != nil {
		return nil, err
	}
	if conversationID == uuid.Nil {
		return nil, nil
	}

	// Step 2: Get conversation info
	var conv entity.Conversation
	if err := r.db.Select("id", "title", "is_group").
		Where("id = ?", conversationID).
		First(&conv).Error; err != nil {
		return nil, err
	}

	// Step 3: Get 20 latest messages with sender info
	var messageResults []dto.LtMessageResponse
	err = r.db.Raw(`
		SELECT m.content, m.sent_at, m.sender_id, u.name AS sender_name
		FROM messages m
		JOIN users u ON u.id = m.sender_id
		WHERE m.conversation_id = ?
		ORDER BY m.sent_at DESC
		LIMIT 20
	`, conversationID).Scan(&messageResults).Error
	if err != nil {
		return nil, err
	}

	// Step 4: Get participants (only id + name)
	var participantResults []dto.LtParticipantResponse
	err = r.db.Raw(`
		SELECT u.id AS user_id, u.name AS user_name
		FROM participants p
		JOIN users u ON u.id = p.user_id
		WHERE p.conversation_id = ?
	`, conversationID).Scan(&participantResults).Error
	if err != nil {
		return nil, err
	}

	return &dto.ConversationDetail{
		ConversationID: conv.ID,
		Title:          conv.Title,
		IsGroup:        conv.IsGroup,
		Messages:       messageResults,
		Participants:   participantResults,
	}, nil
}

func (r *conversationRepository) CheckConversation(ctx context.Context, user1, user2 uuid.UUID, trx *gorm.DB) (int8, error) {
	var conv entity.Conversation

	db := r.db
	if trx != nil {
		db = trx
	}

	subquery := db.
		Table("participants").
		Select("conversation_id").
		Where("user_id IN (?)", []uuid.UUID{user1, user2}).
		Group("conversation_id").
		Having("COUNT(*) = 2")

	err := db.WithContext(ctx).
		Preload("Participants").
		Where("id IN (?) AND is_group = false", subquery).
		First(&conv).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return 1, nil
}

func (r *conversationRepository) FindGroupConversation(ctx context.Context, filter dto.ConversationFilter, userId uuid.UUID) ([]dto.GroupConversation, int64, error) {
	results := []dto.GroupConversation{}
	var total int64

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	listQuery := r.db.
		Table("conversations AS c").
		Select("c.id, c.title, COUNT(cp.user_id) AS total_people").
		Joins("LEFT JOIN participants p ON p.conversation_id = c.id AND p.user_id = ?", userId).
		Joins("JOIN participants cp ON cp.conversation_id = c.id").
		Where("c.is_group = ?", true).
		Where("p.user_id IS NULL"). // user hasn't joined
		Group("c.id, c.title").
		Limit(filter.Limit).
		Offset(filter.Offset)

	countQuery := r.db.
		Table("conversations AS c").
		Select("COUNT(*)").
		Joins("LEFT JOIN participants p ON p.conversation_id = c.id AND p.user_id = ?", userId).
		Where("c.is_group = ?", true).
		Where("p.user_id IS NULL")

	if filter.Title != "" {
		listQuery = listQuery.Where("c.title ILIKE ?", "%"+filter.Title+"%")
		countQuery = countQuery.Where("c.title ILIKE ?", "%"+filter.Title+"%")
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := listQuery.WithContext(ctx).Scan(&results).Error; err != nil {
			errChan <- fmt.Errorf("failed to fetch conversations: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := countQuery.WithContext(ctx).Scan(&total).Error; err != nil {
			errChan <- fmt.Errorf("failed to fetch count: %w", err)
		}
	}()

	wg.Wait()
	close(errChan)

	if len(errChan) > 0 {
		err := <-errChan
		return nil, 0, err
	}

	return results, total, nil
}

func (r *conversationRepository) BeginTrx(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx).Begin()
}
