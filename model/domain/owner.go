package domain

type Owner struct {
	ID     uint64 `json:"id" gorm:"primaryKey"`
	Email  string `json:"email" gorm:"size:255" validate:"email"`
	UserID uint64 `json:"user_id"`
	User   User   `json:"-" gorm:"foreignKey:UserID;references:ID"`
}
