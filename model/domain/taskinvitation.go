package domain

import "time"

type Invitation struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	TaskID    uint64    `json:"task_id"`
	Task      Task      `json:"-" gorm:"foreignKey:TaskID;references:ID"`
	UserID    uint64    `json:"user_id"`
	User      User      `json:"-" gorm:"foreignKey:UserID;references:ID"`
	UserEmail string    `json:"user_email" gorm:"-"` // Add this field
	Role      string    `json:"role" gorm:"type:enum('manager','employee')"`
	Status    string    `json:"status" gorm:"type:enum('pending','accepted','rejected')"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
