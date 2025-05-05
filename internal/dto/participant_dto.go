package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateParticipantRequest struct {
	ConversationID uuid.UUID `json:"conversationId" validate:"required"`
	UserID         uuid.UUID `json:"userId" validate:"required"`
}

type ParticipantResponse struct {
	ConversationID *uuid.UUID `json:"conversationId,omitempty"`
	ID             uuid.UUID  `json:"id"`
	Name           string     `json:"name"`
	JoinedAt       time.Time  `json:"joinedAt"`
}

type ParticipantRawResponse struct {
	ConversationID uuid.UUID `json:"conversationId"`
	UserID         uuid.UUID `json:"userId"`
	JoinedAt       time.Time `json:"joinedAt"`
}
