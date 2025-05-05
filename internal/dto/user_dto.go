package dto

import (
	"time"

	"github.com/google/uuid"
)

type UserFilter struct {
	Limit  int
	Offset int
	UserID uuid.UUID
	Name   string
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email" validate:"omitempty,email"`
}

type UserResponseDetail struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type UserResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
