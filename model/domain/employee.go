package domain

type Employee struct {
	ID               uint64     `json:"id" gorm:"primaryKey"`
	Email            string     `json:"email" gorm:"size:255" validate:"email"`
	UserID           uint64     `json:"user_id"`
	User             User       `json:"-" gorm:"foreignKey:UserID;references:ID"`
	InvitationID     *uint64    `json:"invitation_id,omitempty"`
	InvitationStatus string     `json:"invitation_status,omitempty"`
	Invitation       Invitation `json:"-" gorm:"foreignKey:InvitationID;references:ID"`
	CustomRole       string     `json:"custom_role" gorm:"size:255"`
	TaskID           uint64     `json:"task_id"`
}
