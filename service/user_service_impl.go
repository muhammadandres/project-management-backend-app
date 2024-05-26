package service

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"manajemen_tugas_master/helper"
	"manajemen_tugas_master/model/domain"
	"manajemen_tugas_master/repository"
	"os"
	"time"
)

type userService struct {
	userRepository repository.UserRepository
	validator      *validator.Validate
}

// NewUserService menggabungkan userService dan Userservice untuk membuat instance UserService baru,
// yang memiliki kemampuan UserRepository dan validate
func NewUserService(userRepository repository.UserRepository, validator *validator.Validate) UserService {
	return &userService{userRepository, validator}
}

func (s *userService) SignupUser(user *domain.User) (*domain.User, error) {
	if err := s.validator.Struct(user); err != nil {
		// Jika terjadi kesalahan validasi, konversikan ke satu pesan kesalahan
		var errMsg string
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			errMsg += fmt.Sprintf("Invalid format in %s", fieldError.Field())
		}
		return nil, errors.New(errMsg)
	}

	// Password diubah menjadi hash menggunakan algoritma bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return nil, errors.New("Failed to hash password")
	}

	user.Password = string(hash)

	signup, err := s.userRepository.Signup(user)
	if err != nil {
		return nil, err
	}

	return signup, nil
}

func (s *userService) LoginUser(user *domain.User) (string, error) {
	if err := s.validator.Struct(user); err != nil {
		var errMsg string
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			errMsg += fmt.Sprintf("Invalid format in %s", fieldError.Field())
		}
		return "", errors.New(errMsg)
	}

	// Mendapatkan data pengguna dari repository
	userRepo := *user
	dbUser, err := s.userRepository.Login(&userRepo)
	if err != nil {
		return "", errors.New("User not found") // Mengembalikan pesan kesalahan jika login gagal
	}

	// membandingkan hash password di database, dengan hash password yang baru di kirimkan
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		return "", errors.New("Invalid password") // Mengembalikan pesan kesalahan jika password salah
	}

	// generate token jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": dbUser.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err // Return empty token string and error
	}

	return tokenString, nil // Return token string dan nil error
}

func (s *userService) RequireAuthUser(tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("SECRET")), nil
	})
	if err != nil {
		return nil, err // Mengembalikan error jika terjadi kesalahan saat mem-parse token
	}

	// Memastikan token adalah token yang valid
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token") // Mengembalikan error jika token tidak valid
	}
	if time.Now().Unix() > int64(claims["exp"].(float64)) {
		return nil, errors.New("Token expired") // Mengembalikan error jika token telah kadaluarsa
	}

	// Find the user with token sub
	userID := claims["sub"]
	user, err := s.userRepository.FindById(userID)
	if err != nil {
		return nil, errors.New("User not found")
	}

	return user, nil
}

func (s *userService) GetUserByID(id interface{}) (*domain.User, error) {
	return s.userRepository.FindById(id)
}

func (s *userService) FindAllUsers() ([]*domain.User, error) {
	return s.userRepository.FindAll()
}

func (s *userService) UpdateUser(user *domain.User) (*domain.User, error) {
	if err := s.validator.Struct(user); err != nil {
		// Jika terjadi kesalahan validasi, konversikan ke satu pesan kesalahan
		var errMsg string
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			errMsg += fmt.Sprintf("Invalid format in %s", fieldError.Field())
		}
		// Mengembalikan pesan kesalahan sebagai error
		return nil, errors.New(errMsg)
	}

	updateUser, err := s.userRepository.Update(user)
	if err != nil {
		return nil, errors.New("User not found")
	}

	return updateUser, nil
}

func (s *userService) DeleteUser(id uint) error {
	if id == 0 {
		return errors.New("Invalid user ID")
	}

	db, err := s.userRepository.Delete(id)
	if err != nil {
		return err
	}

	var user *domain.User
	err = helper.ResetAutoIncrement(db, &user, "id", "users")
	if err != nil {
		return err
	}

	return nil
}
