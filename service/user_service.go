package service

import (
	"manajemen_tugas_master/model/domain"
)

type UserService interface {
	SignupUser(user *domain.User) (string, error)
	LoginUser(user *domain.User) (string, error)
	RequireAuthUser(tokenString string) (*domain.User, error)
	GoogleOauth(email string) error
	RequireOauth(email string) (*domain.User, error)
	ForgotPassword(email, newPassword string) error
	GetUserByID(id interface{}) (*domain.User, error)
	FindAllUsers() ([]*domain.User, error)
	UpdateUser(user *domain.User) (*domain.User, error)
	DeleteUser(id uint) error
}
