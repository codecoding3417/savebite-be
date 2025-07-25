package dto

import "github.com/google/uuid"

type CreateUserRequest struct {
	Name  string
	Email string
}

type UserParam struct {
	ID    uuid.UUID
	Name  string
	Email string
}
