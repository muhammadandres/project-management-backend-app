package repository

import (
	"errors"
	"fmt"
	"manajemen_tugas_master/model/domain"

	"gorm.io/gorm"
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

func (r *userRepository) Signup(user *domain.User) error {
	var count int64
	r.db.First(&user, "email = ?", user.Email).Count(&count)
	if count == 0 {
		if err := r.db.Create(&user).Error; err != nil {
			return fmt.Errorf("err %v", err)
		}
	}
	if count > 0 {
		return errors.New("email already exists, please login instead")
	}

	return nil
}

func (r *userRepository) Login(user *domain.User) (*domain.User, error) {
	var dbUser domain.User
	if err := r.db.First(&dbUser, "email = ?", user.Email).Error; err != nil {
		return nil, err
	}
	return &dbUser, nil
}

func (r *userRepository) GoogleOauth(email string) error {
	var user domain.User
	var count int64

	// Periksa apakah user dengan email tersebut sudah ada
	if err := r.db.Model(&domain.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		// Jika user belum ada, buat user baru dengan password googleauth
		user.Email = email
		user.Password = "GoogleAuth"
		if err := r.db.Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

func (r *userRepository) RequireOauth(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("User not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdatePassword(userID uint64, newPassword string) error {
	result := r.db.Model(&domain.User{}).Where("id = ?", userID).Update("password", newPassword)
	return result.Error
}

func (r *userRepository) FindById(id interface{}) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
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
