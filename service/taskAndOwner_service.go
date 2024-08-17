package service

import (
	"manajemen_tugas_master/model/domain"
	"manajemen_tugas_master/model/web"
)

type TaskAndOwnerService interface {
	CreateTaskAndOwner(user *domain.User, task *domain.Task, board *domain.Board) (*domain.Task, *domain.Owner, error)
	GetTaskAndOwnerById(id uint) (*domain.TaskWithInvitation, error)
	FindAllTasksAndOwners() ([]*domain.TaskWithInvitation, error)
	FindAllOwners() ([]*domain.Task, error)
	FindAllManagers() ([]*domain.Task, error)
	FindAllEmployees() ([]*domain.Task, error)
	FindAllPlanningFiles() ([]*domain.Task, error)
	FindAllProjectFiles() ([]*domain.Task, error)
	UpdateTaskAndOwner(task *domain.Task, newManagers []domain.Manager, newEmployees []domain.Employee, planningFile *domain.PlanningFile, projectFile *domain.ProjectFile, taskID uint, boardID uint) (*web.UpdateResponse, error)
	UpdateValidationOwner(taskID uint, userID uint) error
	UpdateValidationManager(taskID uint, userID uint) error
	UpdateValidationEmployee(taskID uint, userID uint) error
	DeleteManager(taskId uint, managerId uint) error
	DeleteEmployee(taskId uint, employeeId uint) error
	DeletePlanningFile(fileId uint) (string, error)
	DeleteProjectFile(fileId uint) (string, error)
	DeleteTaskAndOwner(taskID uint) error

	UpdateOwnerCustomRole(taskID uint, customRole string) (*domain.Owner, error)
	AddManager(taskID uint, email string) (*web.ManagerResponse, error)
	AddEmployee(taskID uint, email string) (*web.EmployeeResponse, error)
	UpdateManagerEmail(taskID uint, oldEmail, newEmail string) (*web.ManagerResponse, error)
	UpdateEmployeeEmail(taskID uint, oldEmail, newEmail string) (*web.EmployeeResponse, error)
	UpdateManagerCustomRole(taskID uint, email, customRole string) error
	UpdateEmployeeCustomRole(taskID uint, email, customRole string) error

	GetManagersByTaskID(taskID uint) ([]web.ManagerResponse, error)
	GetEmployeesByTaskID(taskID uint) ([]web.EmployeeResponse, error)

	GetAllInvitations() ([]domain.Invitation, error)
	RespondToInvitation(invitationID uint64, response string, role string) (*domain.Invitation, error)
}
