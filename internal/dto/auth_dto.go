package dto

import "github.com/google/uuid"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterRequest struct {
	Name     string `json:"name"  validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	UserId     uuid.UUID `json:"userId"`
	UserName   string    `json:"userName"`
	AccesToken string    `json:"accessToken"`
}
