package domain

import "time"

type PlanningFile struct {
	ID        uint64    `json:"id" gorm:"primaryKey"`
	FileUrl   string    `json:"file_url" gorm:"size:255"`
	FileName  string    `json:"file_name" gorm:"size:255"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}
