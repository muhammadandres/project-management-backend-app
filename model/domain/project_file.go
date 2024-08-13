package domain

type ProjectFile struct {
	ID       uint64 `json:"id" gorm:"primaryKey"`
	FileUrl  string `json:"fileUrl" gorm:"size:255"`
	FileName string `json:"FileName" gorm:"size:255"`
}
