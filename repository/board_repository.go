package repository

import (
	"manajemen_tugas_master/model/domain"

	"gorm.io/gorm"
)

type BoardRepository interface {
	CreateBoard(board *domain.Board) (*domain.Board, error)
	FindById(id uint64) (*domain.Board, error)
	GetAllBoards() ([]*domain.Board, error)
	DeleteById(id uint64) (*gorm.DB, int64, int64, int64, int64, int64, error)
	EditBoard(id uint64, newNameBoard string) (*domain.Board, error)
}
