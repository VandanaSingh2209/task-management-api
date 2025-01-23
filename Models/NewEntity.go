package Models

import "time"

type AmplTaskList1 struct {
	ID          uint      `gorm:"column:ID" json:"ID"`
	Title       string    `gorm:"not null" json:"title" binding:"required"`
	Description string    `gorm:"not null" json:"description" binding:"required"`
	Status      string    `gorm:"default:'pending'" json:"status" binding:"required,oneof=pending in-progress completed"`
	CreatedAt   time.Time `gorm:"column:CreatedAt"`
	UpdatedAt   time.Time `gorm:"column:UpdatedAt"`
}

func (b *AmplTaskList1) TableName() string {
	return "AmplTaskList"
}
