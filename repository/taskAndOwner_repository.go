package repository

import (
	"manajemen_tugas_master/model/domain"

	"gorm.io/gorm"
)

type TaskAndOwnerRepository interface {
	Create(user *domain.User, task *domain.Task, board *domain.Board) (*domain.Task, *domain.Owner, error)
	FindById(id uint) (*domain.TaskWithInvitation, error)
	FindAll() ([]*domain.TaskWithInvitation, error)
	FindAllOwners() ([]*domain.Task, error)
	FindAllManagers() ([]*domain.Task, error)
	FindAllEmployees() ([]*domain.Task, error)
	FindAllPlanningFiles() ([]*domain.Task, error)
	FindAllProjectFiles() ([]*domain.Task, error)
	GetNameEmailsDescription(taskID uint64) (ownerEmail string, managerEmails []string, employeeEmails []string, nametask string, description string, err error)
	Update(task *domain.Task, planningFile *domain.PlanningFile, projectFile *domain.ProjectFile) (*domain.Task, *domain.PlanningFile, *domain.ProjectFile, error)
	UpdateValidationOwner(taskID uint, userID uint) error
	UpdateValidationManager(taskID uint, userID uint) error
	UpdateValidationEmployee(taskID uint, userID uint) error
	DeleteManager(taskId uint, managerId uint) (*gorm.DB, int64, int64, int64, error)
	DeleteEmployee(taskId uint, employeeId uint) (*gorm.DB, int64, error)
	DeletePlanningFile(fileId uint) (*gorm.DB, string, error)
	DeleteProjectFile(fileId uint) (*gorm.DB, string, error)
	Delete(taskID uint) (*gorm.DB, int64, int64, int64, int64, int64, error)

	CreateInvitation(invitation *domain.Invitation) (*domain.Invitation, error)
	FindInvitationByID(id uint64) (*domain.Invitation, error)
	UpdateInvitation(invitation *domain.Invitation) error
	GetAllInvitations() ([]domain.Invitation, error)
	UpdateManagerInvitationStatus(invitationID uint64, status string) (*domain.Invitation, error)
	UpdateEmployeeInvitationStatus(invitationID uint64, status string) (*domain.Invitation, error)

	UpdateOwnerCustomRole(taskID uint, customRole string) (*domain.Owner, error)
	AddManager(taskID uint, email string) (*domain.Manager, error)
	AddEmployee(taskID uint, email string) (*domain.Employee, error)
	UpdateManagerEmail(taskID uint, oldEmail, newEmail string) (*domain.Manager, error)
	UpdateEmployeeEmail(taskID uint, oldEmail, newEmail string) (*domain.Employee, error)
	UpdateManagerCustomRole(taskID uint, email, customRole string) error
	UpdateEmployeeCustomRole(taskID uint, email, customRole string) error
	GetManagersByTaskID(taskID uint) ([]domain.Manager, error)
	GetEmployeesByTaskID(taskID uint) ([]domain.Employee, error)
}
