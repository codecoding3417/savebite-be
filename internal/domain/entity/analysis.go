package entity

import (
	"github.com/google/uuid"
	"time"
)

type Analysis struct {
	ID          uuid.UUID `gorm:"type:char(36);primaryKey"`
	UserID      uuid.UUID `gorm:"not null"`
	Feedback    string    `gorm:"type:text;not null"`
	Ingredients []Ingredient
	CreatedAt   time.Time `gorm:"type:timestamp;autoCreateTime"`
}
