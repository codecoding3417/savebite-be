package entity

import (
	"github.com/google/uuid"
	"time"
)

type Ingredient struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	AnalysisID uuid.UUID `gorm:"not null;"`
	Name       string    `gorm:"size:100;not null"`
	Status     string    `gorm:"type:enum('useable','unuseable');not null"`
	CreatedAt  time.Time `gorm:"type:timestamp;autoCreateTime"`
}
