package entity

import (
	"time"

	"github.com/google/uuid"
)

type Conversation struct {
	ID           uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Title        string        `gorm:"not null"`
	IsGroup      bool          `gorm:"default:false"`
	CreatedAt    time.Time     `gorm:"autoCreateTime"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime"`
	Participants []Participant `gorm:"foreignKey:ConversationID"`
}

type ConversationSummary struct {
	ConversationID  uuid.UUID
	Title           string
	IsGroup         bool
	LatestMessage   string
	MessageSentAt   time.Time
	ParticipantName string
}
