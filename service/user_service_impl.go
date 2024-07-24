package service

import (
	"errors"
	"fmt"
	"manajemen_tugas_master/helper"
	"manajemen_tugas_master/model/domain"
	"manajemen_tugas_master/repository"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

func (s *userService) SignupUser(user *domain.User) (string, error) {
	if err := s.validator.Struct(user); err != nil {
		// Jika terjadi kesalahan validasi, konversikan ke satu pesan kesalahan
		var errMsg string
		validationErrors := err.(validator.ValidationErrors)
		for _, fieldError := range validationErrors {
			errMsg += fmt.Sprintf("Invalid format in %s", fieldError.Field())
		}
		return "", errors.New(errMsg)
	}

	// Password diubah menjadi hash menggunakan algoritma bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return "", errors.New("Failed to hash password")
	}

	user.Password = string(hash)

	if err := s.userRepository.Signup(user); err != nil {
		return "", err
	}

	// ambil id user berdasarkan email
	dbUser, err := s.userRepository.Login(user)
	if err != nil {
		return "", errors.New("User not found") // Mengembalikan pesan kesalahan jika login gagal
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

	return tokenString, nil
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

	if user.Password == "" {
		return "", errors.New("Password is required")
	}

	// Simpan password yang diberikan user
	providedPassword := user.Password

	// Mendapatkan data pengguna dari repository
	dbUser, err := s.userRepository.Login(user)
	if err != nil {
		return "", errors.New("User not found")
	}

	// membandingkan hash password di database, dengan password yang baru dikirimkan
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(providedPassword))
	if err != nil {
		return "", errors.New("Invalid password")
	}

	// generate token jwt
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": dbUser.ID,
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
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
		return nil, err
	}

	return user, nil
}

func (s *userService) GoogleOauth(email string) error {
	return s.userRepository.GoogleOauth(email)
}

func (s *userService) RequireOauth(email string) (*domain.User, error) {
	user, err := s.userRepository.RequireOauth(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) ForgotPassword(email, newPassword string) error {
	if err := s.validator.Var(email, "required,email"); err != nil {
		return errors.New("Invalid format in Email")
	}

	if err := s.validator.Var(newPassword, "required,min=6"); err != nil {
		return errors.New("Invalid password format")
	}

	user, err := s.userRepository.GetUserByEmail(email)
	if err != nil {
		return errors.New("User not found")
	}

	// Hash kata sandi baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 10)
	if err != nil {
		return errors.New("Failed to hash new password")
	}

	// Simpan kata sandi baru yang sudah di-hash
	err = s.userRepository.UpdatePassword(user.ID, string(hashedPassword))
	if err != nil {
		return errors.New("Failed to update password")
	}

	return nil
}

func (s *userService) GetUserByID(id interface{}) (*domain.User, error) {
	return s.userRepository.FindById(id)
}

func (s *userService) FindAllUsers() ([]*domain.User, error) {
	// test google calendar
	// event, err := helper.CreateGoogleCalendarEvent(
	// 	userEmail,                        // user email
	// 	"New Task Due Date",              // summary
	// 	"Task description",               // description
	// 	"2024-06-15T09:00:00Z",           // startDateTime
	// 	"2024-06-15T17:00:00Z",           // endDateTime
	// 	"UTC",                            // timeZone
	// 	[]string{"attendee@example.com"}, // attendees
	// )
	// if err != nil {
	// 	log.Printf("Error creating Google Calendar event: %v", err)
	// } else {
	// 	log.Printf("Event created: %s", event.HtmlLink)
	// }

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
