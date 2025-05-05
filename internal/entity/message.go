package entity

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ConversationID uuid.UUID `gorm:"type:uuid;not null"`
	SenderID       uuid.UUID `gorm:"type:uuid;not null"`
	Content        string    `gorm:"not null"`
	SentAt         time.Time `gorm:"autoCreateTime"`

	User User `gorm:"foreignKey:SenderID"`
}
