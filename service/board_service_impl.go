package service

import (
	"errors"
	"manajemen_tugas_master/helper"
	"manajemen_tugas_master/model/domain"
	"manajemen_tugas_master/repository"
)

type boardService struct {
	boardRepository repository.BoardRepository
}

func NewBoardService(boardRepository repository.BoardRepository) BoardService {
	return &boardService{boardRepository}
}

func (s *boardService) CreateBoard(board *domain.Board) (*domain.Board, error) {
	if board.NameBoard == "" {
		return nil, errors.New("Masukkan nama board terlebih dahulu")
	}

	if board.UserID == 0 {
		return nil, errors.New("User ID tidak valid")
	}

	boardDb, err := s.boardRepository.CreateBoard(board)
	if err != nil {
		return nil, err
	}

	return boardDb, nil
}

func (s *boardService) GetBoardById(id uint64) (*domain.Board, error) {
	board, err := s.boardRepository.FindById(id)
	if err != nil {
		return nil, err
	}
	return board, nil
}

func (s *boardService) GetAllBoards() ([]*domain.Board, error) {
	return s.boardRepository.GetAllBoards()
}

func (s *boardService) DeleteBoardById(id uint64) error {
	db, countTasks, countManagers, countEmployees, countPlanningFiles, countProjectFiles, err := s.boardRepository.DeleteById(id)
	if err != nil {
		return err
	}

	// Reset auto increment
	if countTasks > 0 {
		var task domain.Task
		err = helper.ResetAutoIncrement(db, &task, "id", "tasks")
		if err != nil {
			return err
		}
	}

	if countManagers > 0 {
		var manager domain.Manager
		err = helper.ResetAutoIncrement(db, &manager, "id", "managers")
		if err != nil {
			return err
		}
	}

	if countEmployees > 0 {
		var employee domain.Employee
		err = helper.ResetAutoIncrement(db, &employee, "id", "employees")
		if err != nil {
			return err
		}
	}

	if countPlanningFiles > 0 {
		var planningFile domain.PlanningFile
		err = helper.ResetAutoIncrement(db, &planningFile, "id", "planning_files")
		if err != nil {
			return err
		}
	}

	if countProjectFiles > 0 {
		var projectFile domain.ProjectFile
		err = helper.ResetAutoIncrement(db, &projectFile, "id", "project_files")
		if err != nil {
			return err
		}
	}

	var board domain.Board
	err = helper.ResetAutoIncrement(db, &board, "id", "boards")
	if err != nil {
		return err
	}

	return nil
}

func (s *boardService) EditBoard(id uint64, newNameBoard string) (*domain.Board, error) {
	if newNameBoard == "" {
		return nil, errors.New("New board name cannot be empty")
	}

	return s.boardRepository.EditBoard(id, newNameBoard)
}
