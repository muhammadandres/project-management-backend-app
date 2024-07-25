package domain

import "time"

type Board struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	NameBoard string    `json:"name_board" gorm:"size:255"`
	Tasks     []Task    `json:"tasks" gorm:"foreignKey:BoardID;references:ID"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}