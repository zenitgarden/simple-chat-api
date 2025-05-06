package dto

import (
	"time"

	"github.com/google/uuid"
)

type ConversationFilter struct {
	Title  string
	Limit  int
	Offset int
	Sort   string
}

type CreateConversationRequest struct {
	Title   string    `json:"title" validate:"required"`
	IsGroup bool      `json:"isGroup"`
	UserId  uuid.UUID `json:"userId"`
}

type UpdateConversationRequest struct {
	Title   string `json:"title"`
	IsGroup *bool  `json:"isGroup"`
}

type ConversationResponse struct {
	ID           uuid.UUID              `json:"id"`
	Title        string                 `json:"title"`
	IsGroup      bool                   `json:"isGroup"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
	Participants *[]ParticipantResponse `json:"participants,omitempty"`
}

type AllConversationResponse struct {
	ID          uuid.UUID  `json:"id"`
	Title       string     `json:"title"`
	IsGroup     bool       `json:"isGroup"`
	LastMessage string     `json:"lastMessage"`
	SentAt      *time.Time `json:"sentAt"`
	CreatedBy   uuid.UUID  `json:"createdBy"`
}

type LtMessageResponse struct {
	Content    string    `json:"content"`
	SentAt     time.Time `json:"sent_at"`
	SenderID   uuid.UUID `json:"sender_id"`
	SenderName string    `json:"sender_name"`
}

// dto/participant_response.go
type LtParticipantResponse struct {
	UserID   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`
}

// dto/conversation_detail.go
type ConversationDetail struct {
	ConversationID uuid.UUID               `json:"conversation_id"`
	Title          string                  `json:"title"`
	IsGroup        bool                    `json:"is_group"`
	Messages       []LtMessageResponse     `json:"messages"`
	Participants   []LtParticipantResponse `json:"participants"`
}
