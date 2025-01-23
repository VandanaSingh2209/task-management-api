package Models

import (
	"time"
)

type AmplTaskList struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Title       string `gorm:"not null" json:"title" binding:"required"`
	Description string `gorm:"not null" json:"description" binding:"required"`
	Status      string `gorm:"default:'pending'" json:"status" binding:"required,oneof=pending in-progress completed"`
	// CreatedAt   time.Time `json:"created_at"`
	// UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt time.Time `gorm:"column:CreatedAt"`
	UpdatedAt time.Time `gorm:"column:UpdatedAt"`
}

func (b *AmplTaskList) TableName() string {
	return "AmplTaskList"
}
