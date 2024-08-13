package domain

type Board struct {
	ID        uint64 `json:"id" gorm:"primaryKey"`
	NameBoard string `json:"name_board" gorm:"size:255"`
	UserID    uint64 `json:"user_id"`
	User      User   `json:"user" gorm:"foreignKey:UserID;references:ID"`
	Tasks     []Task `json:"tasks" gorm:"foreignKey:BoardID;references:ID"`
}
