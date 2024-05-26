package repository

import (
	"gorm.io/gorm"
	"manajemen_tugas_master/model/domain"
)

// UserRepository adalah interface untuk operasi-operasi yang berhubungan dengan entitas User
type UserRepository interface {
	Signup(user *domain.User) (*domain.User, error)
	Login(user *domain.User) (*domain.User, error)
	FindById(id interface{}) (*domain.User, error)
	FindAll() ([]*domain.User, error)
	Update(user *domain.User) (*domain.User, error)
	Delete(id uint) (*gorm.DB, error)
}
