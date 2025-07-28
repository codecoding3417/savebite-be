package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name      string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"type:varchar(100);not null"`
	Analyses  []Analysis
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
}
