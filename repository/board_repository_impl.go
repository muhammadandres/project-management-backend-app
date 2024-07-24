package repository

import (
	"errors"
	"fmt"
	"manajemen_tugas_master/model/domain"

	"gorm.io/gorm"
)

type boardRepository struct {
	db *gorm.DB
}

func NewBoardRepository(db *gorm.DB) BoardRepository {
	return &boardRepository{db}
}

func (b *boardRepository) CreateBoard(board *domain.Board) (*domain.Board, error) {
	boardDb := &domain.Board{
		NameBoard: board.NameBoard,
	}

	if err := b.db.Create(boardDb).Error; err != nil {
		return nil, fmt.Errorf("err %v", err)
	}

	return boardDb, nil
}

func (b *boardRepository) FindById(id uint64) (*domain.Board, error) {
	var board domain.Board
	if err := b.db.Preload("Tasks").Preload("Tasks.Owner").Preload("Tasks.Manager").Preload("Tasks.Employee").Preload("Tasks.PlanningFile").Preload("Tasks.ProjectFile").First(&board, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("board not found")
		}
		return nil, err
	}
	return &board, nil
}

func (b *boardRepository) GetAllBoards() ([]*domain.Board, error) {
	var boards []*domain.Board
	if err := b.db.Preload("Tasks").
		Preload("Tasks.Owner").
		Preload("Tasks.Manager").
		Preload("Tasks.Employee").
		Preload("Tasks.PlanningFile").
		Preload("Tasks.ProjectFile").
		Find(&boards).Error; err != nil {
		return nil, err
	}
	return boards, nil
}

func (b *boardRepository) DeleteById(id uint64) (*gorm.DB, int64, int64, int64, int64, int64, error) {
	var (
		countTasks         int64
		countManagers      int64
		countEmployees     int64
		countPlanningFiles int64
		countProjectFiles  int64
	)

	err := b.db.Transaction(func(tx *gorm.DB) error {
		// Fetch all task IDs associated with this board
		var taskIDs []uint
		if err := tx.Model(&domain.Task{}).Where("board_id = ?", id).Pluck("id", &taskIDs).Error; err != nil {
			return err
		}

		// Count and delete associated records in task_managers
		tx.Model(&domain.Manager{}).Where("task_id IN (?)", taskIDs).Count(&countManagers)
		if err := tx.Exec("DELETE FROM task_managers WHERE task_id IN (?)", taskIDs).Error; err != nil {
			return err
		}

		// Count and delete associated records in task_employees
		tx.Model(&domain.Employee{}).Where("task_id IN (?)", taskIDs).Count(&countEmployees)
		if err := tx.Exec("DELETE FROM task_employees WHERE task_id IN (?)", taskIDs).Error; err != nil {
			return err
		}

		// Count and delete associated records in task_planning_files
		tx.Model(&domain.PlanningFile{}).Where("task_id IN (?)", taskIDs).Count(&countPlanningFiles)
		if err := tx.Exec("DELETE FROM task_planning_files WHERE task_id IN (?)", taskIDs).Error; err != nil {
			return err
		}

		// Count and delete associated records in task_project_files
		tx.Model(&domain.ProjectFile{}).Where("task_id IN (?)", taskIDs).Count(&countProjectFiles)
		if err := tx.Exec("DELETE FROM task_project_files WHERE task_id IN (?)", taskIDs).Error; err != nil {
			return err
		}

		// Count and delete associated tasks
		tx.Model(&domain.Task{}).Where("board_id = ?", id).Count(&countTasks)
		if err := tx.Where("board_id = ?", id).Delete(&domain.Task{}).Error; err != nil {
			return err
		}

		// Delete the board
		if err := tx.Delete(&domain.Board{}, id).Error; err != nil {
			return err
		}

		return nil
	})

	return b.db, countTasks, countManagers, countEmployees, countPlanningFiles, countProjectFiles, err
}

func (b *boardRepository) EditBoard(id uint64, newNameBoard string) (*domain.Board, error) {
	var board domain.Board

	// Check if the board exists
	if err := b.db.First(&board, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("board not found")
		}
		return nil, err
	}

	// Update the board name
	board.NameBoard = newNameBoard
	if err := b.db.Save(&board).Error; err != nil {
		return nil, err
	}

	return &board, nil
}
