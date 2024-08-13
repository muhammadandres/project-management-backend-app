package domain

type User struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	Email    string `json:"email" gorm:"size:255" validate:"email"`
	Password string `json:"password" gorm:"size:255"`
}
