package service

import (
	"manajemen_tugas_master/model/domain"
	"manajemen_tugas_master/model/web"
)

type TaskAndOwnerService interface {
	CreateTaskAndOwner(user *domain.User, task *domain.Task, board *domain.Board) (*domain.Task, *domain.Owner, error)
	GetTaskAndOwnerById(id uint) (*domain.Task, error)
	FindAllTasksAndOwners() ([]*domain.Task, error)
	FindAllOwners() ([]*domain.Task, error)
	FindAllManagers() ([]*domain.Task, error)
	FindAllEmployees() ([]*domain.Task, error)
	FindAllPlanningFiles() ([]*domain.Task, error)
	FindAllProjectFiles() ([]*domain.Task, error)
	UpdateTaskAndOwner(task *domain.Task, manager *domain.Manager, employee *domain.Employee, planningFile *domain.PlanningFile, projectFile *domain.ProjectFile, taskID uint, boardID uint) (*web.UpdateResponse, error)
	UpdateValidationOwner(taskID uint, userID uint) error
	UpdateValidationManager(taskID uint, userID uint) error
	UpdateValidationEmployee(taskID uint, userID uint) error
	DeleteManager(taskId uint, managerId uint) error
	DeleteEmployee(taskId uint, employeeId uint) error
	DeletePlanningFile(fileId uint) (string, error)
	DeleteProjectFile(fileId uint) (string, error)
	DeleteTaskAndOwner(taskID uint) error
}
