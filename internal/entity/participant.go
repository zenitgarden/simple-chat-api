package entity

import (
	"time"

	"github.com/google/uuid"
)

type Participant struct {
	ID             uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	ConversationID uuid.UUID `gorm:"type:uuid;not null"`
	UserID         uuid.UUID `gorm:"type:uuid;not null"`
	JoinedAt       time.Time `gorm:"autoCreateTime"`
	User           User      `gorm:"foreignKey:UserID"`
}
