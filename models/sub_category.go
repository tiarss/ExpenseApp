package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubCategory struct {
	gorm.Model
	ID         uuid.UUID `gorm:"type:char(36);primaryKey;" json:"id"`
	Name       string    `json:"name"`
	CategoryID uuid.UUID `json:"category_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (s *SubCategory) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = uuid.New()
	return
}
