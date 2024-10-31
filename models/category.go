package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	ID            uuid.UUID     `gorm:"type:char(36);primaryKey;" json:"id"`
	Name          string        `json:"name"`
	CategoryType  string        `json:"category_type"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	SubCategories []SubCategory `gorm:"foreignKey:CategoryID;references:ID" json:"sub_categories"`
}

func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}
