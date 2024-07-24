package service

import "manajemen_tugas_master/model/domain"

type BoardService interface {
	CreateBoard(board *domain.Board) (*domain.Board, error)
	GetBoardById(id uint64) (*domain.Board, error)
	GetAllBoards() ([]*domain.Board, error)
	DeleteBoardById(id uint64) error
	EditBoard(id uint64, newNameBoard string) (*domain.Board, error)
}
