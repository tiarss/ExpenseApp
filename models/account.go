package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID        uuid.UUID `gorm:"type:char(36);primaryKey;" json:"id"`
	Name      string    `json:"name"`
	UserID    uuid.UUID `json:"user_id"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = uuid.New()
	return
}
