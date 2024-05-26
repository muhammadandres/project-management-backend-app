package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"manajemen_tugas_master/model/domain"
)

// userRepository adalah implementasi dari UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository menggabungkan userRepository dan UserRepository untuk membuat instance UserRepository baru,
// yang memiliki kemampuan gorm
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Signup(user *domain.User) (*domain.User, error) {
	if err := r.db.Create(&user).Error; err != nil {
		return nil, errors.New("Duplicate entry for user")
	}
	return user, nil
}

func (r *userRepository) Login(user *domain.User) (*domain.User, error) {
	if err := r.db.First(&user, "email = ?", user.Email).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindById(id interface{}) (*domain.User, error) {
	var user *domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) FindAll() ([]*domain.User, error) {
	var users []*domain.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Update(user *domain.User) (*domain.User, error) {
	if err := r.db.Save(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepository) Delete(id uint) (*gorm.DB, error) {
	var user *domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %v", err)
	}

	if err := r.db.Delete(&user).Error; err != nil {
		return nil, fmt.Errorf("failed to delete related users, because they are associated with tasks")
	}

	return r.db, nil
}
