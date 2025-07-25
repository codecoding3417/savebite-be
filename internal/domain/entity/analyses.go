package entity

import "github.com/google/uuid"

type Analyses struct {
	ID uuid.UUID `gorm:"type:char(36);primaryKey"`
}
