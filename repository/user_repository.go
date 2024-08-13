package repository

import (
	"manajemen_tugas_master/model/domain"

	"gorm.io/gorm"
)

// UserRepository adalah interface untuk operasi-operasi yang berhubungan dengan entitas User
type UserRepository interface {
	Signup(user *domain.User) error
	Login(user *domain.User) (*domain.User, error)
	GoogleOauth(email string) error
	RequireOauth(email string) (*domain.User, error)
	GetUserByEmail(email string) (*domain.User, error)
	UpdatePassword(userID uint64, newPassword string) error
	FindById(id interface{}) (*domain.User, error)
	FindAll() ([]*domain.User, error)
	Update(user *domain.User) (*domain.User, error)
	Delete(id uint) (*gorm.DB, error)
}
