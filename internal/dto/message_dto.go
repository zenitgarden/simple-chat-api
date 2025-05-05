package dto

import (
	"time"

	"github.com/google/uuid"
)

type MessageFilter struct {
	Limit          int
	Offset         int
	ConversationID uuid.UUID
}

type CreateMessageRequest struct {
	Type           string    `json:"type"`
	ConversationID uuid.UUID `json:"conversationId"`
	Content        string    `json:"content"`
	SenderID       uuid.UUID `json:"senderId"`
}

type MessageResponse struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversationId"`
	SenderID       uuid.UUID `json:"senderId"`
	SenderName     string    `json:"senderName"`
	Content        string    `json:"content"`
	SentAt         time.Time `json:"sentAt"`
}
