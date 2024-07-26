package service

import (
	"errors"
	"manajemen_tugas_master/helper"
	"manajemen_tugas_master/model/domain"
	"manajemen_tugas_master/model/web"
	"manajemen_tugas_master/repository"

	"github.com/go-playground/validator/v10"
)

type taskAndOwnerService struct {
	taskAndOwnerRepository repository.TaskAndOwnerRepository
	boardRepository        repository.BoardRepository
	validator              *validator.Validate
}

func NewTaskAndOwnerService(taskAndOwnerRepository repository.TaskAndOwnerRepository, boardRepository repository.BoardRepository, validator *validator.Validate) TaskAndOwnerService {
	return &taskAndOwnerService{taskAndOwnerRepository, boardRepository, validator}
}

func (t *taskAndOwnerService) CreateTaskAndOwner(user *domain.User, task *domain.Task, board *domain.Board) (*domain.Task, *domain.Owner, error) {
	if task.NameTask == "" {
		return nil, nil, errors.New("Masukkan nama task terlebih dahulu")
	}

	if board.ID == 0 {
		return nil, nil, errors.New("latestBoardID tidak ditemukan")
	}

	taskDB, ownerDB, err := t.taskAndOwnerRepository.Create(user, task, board)
	if err != nil {
		return nil, nil, err
	}

	return taskDB, ownerDB, nil
}

func (t *taskAndOwnerService) GetTaskAndOwnerById(id uint) (*domain.Task, error) {
	return t.taskAndOwnerRepository.FindById(id)
}

func (t *taskAndOwnerService) FindAllTasksAndOwners() ([]*domain.Task, error) {
	return t.taskAndOwnerRepository.FindAll()
}

func (t *taskAndOwnerService) FindAllOwners() ([]*domain.Task, error) {
	return t.taskAndOwnerRepository.FindAllOwners()
}

func (t *taskAndOwnerService) FindAllManagers() ([]*domain.Task, error) {
	return t.taskAndOwnerRepository.FindAllManagers()
}

func (t *taskAndOwnerService) FindAllEmployees() ([]*domain.Task, error) {
	return t.taskAndOwnerRepository.FindAllEmployees()
}

func (t *taskAndOwnerService) FindAllPlanningFiles() ([]*domain.Task, error) {
	return t.taskAndOwnerRepository.FindAllPlanningFiles()
}

func (t *taskAndOwnerService) FindAllProjectFiles() ([]*domain.Task, error) {
	return t.taskAndOwnerRepository.FindAllProjectFiles()
}

func (t *taskAndOwnerService) UpdateTaskAndOwner(task *domain.Task, manager *domain.Manager, employee *domain.Employee, planningFile *domain.PlanningFile, projectFile *domain.ProjectFile, taskID uint, boardID uint) (*web.UpdateResponse, error) {
	boardDB, err := t.boardRepository.FindById(uint64(boardID))
	if err != nil {
		return nil, err
	}
	// Update task dengan data dari database
	task.BoardID = boardDB.ID

	taskDB, err := t.taskAndOwnerRepository.FindById(taskID)
	if err != nil {
		return nil, err
	}
	// Update task dengan data dari database
	task.ID = taskDB.ID
	task.OwnerID = taskDB.OwnerID

	updateTask, updateManager, updateEmployee, updatePlanningFile, updateProjectFile, err := t.taskAndOwnerRepository.Update(task, manager, employee, planningFile, projectFile)
	if err != nil {
		return nil, err
	}

	// Persiapan respons
	response := &web.UpdateResponse{}

	// Populate response dengan data dari updateTask
	response.NameTask = updateTask.NameTask
	response.PlanningDescription = updateTask.PlanningDescription
	response.PlanningStatus = updateTask.PlanningStatus
	// Send email if planning status is "Approved"
	// if updateTask.PlanningStatus == "Approved" {
	// 	ownerEmail, managerEmails, employeeEmails, err := t.taskAndOwnerRepository.GetTaskEmails(uint64(taskID))
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	to := []string{ownerEmail}
	// 	to = append(to, managerEmails...)
	// 	to = append(to, employeeEmails...)

	// 	subject := "Task Planning Approved"
	// 	body := fmt.Sprintf("The planning for task '%s' has been approved.", updateTask.NameTask)

	// 	err = helper.SendEmail(to, subject, body)
	// 	if err != nil {
	// 		log.Printf("Failed to send email: %v", err)
	// 		return nil, fmt.Errorf("failed to send email: %v", err)
	// 	}
	// }
	// if updateTask.PlanningStatus == "Approved" {
	// 	// Menggunakan alamat email statis untuk pengujian
	// 	to := []string{"riyannasgianjar0401@gmail.com"}

	// 	subject := "Task Planning Approved"
	// 	body := fmt.Sprintf("The planning for task has been approved.")

	// 	log.Printf("Attempting to send email to: %v", to)
	// 	log.Printf("Subject: %s", subject)
	// 	log.Printf("Body: %s", body)

	// 	err = helper.SendEmail(to, subject, body)
	// 	if err != nil {
	// 		log.Printf("Failed to send email: %v", err)
	// 		return nil, fmt.Errorf("failed to send email: %v", err)
	// 	}
	// 	log.Println("Email sent successfully")
	// }

	response.ProjectStatus = updateTask.ProjectStatus
	response.PlanningDueDate = updateTask.PlanningDueDate
	response.ProjectDueDate = updateTask.ProjectDueDate
	response.Priority = updateTask.Priority
	response.ProjectComment = updateTask.ProjectComment

	// Populate managerResponse dengan data dari updateManager jika tidak kosong
	if updateManager.ID != 0 || updateManager.Email != "" || updateManager.UserID != 0 {
		response.Manager.ID = updateManager.ID
		response.Manager.Email = updateManager.Email
		response.Manager.UserID = updateManager.UserID
	}

	// Populate employeeResponse dengan data dari updateEmployee jika tidak kosong
	if updateEmployee.ID != 0 || updateEmployee.Email != "" || updateEmployee.UserID != 0 {
		response.Employee.ID = updateEmployee.ID
		response.Employee.Email = updateEmployee.Email
		response.Employee.UserID = updateEmployee.UserID
	}

	// Populate planningFileResponse dengan data dari updatePlanningFile jika tidak kosong
	if updatePlanningFile.ID != 0 || updatePlanningFile.FileUrl != "" || updatePlanningFile.FileName != "" {
		response.PlanningFile.ID = updatePlanningFile.ID
		response.PlanningFile.FileUrl = updatePlanningFile.FileUrl
		response.PlanningFile.FileName = updatePlanningFile.FileName
	}

	// Populate projectFileResponse dengan data dari updateProjectFile jika tidak kosong
	if updateProjectFile.ID != 0 || updateProjectFile.FileUrl != "" || updateProjectFile.FileName != "" {
		response.ProjectFile.ID = updateProjectFile.ID
		response.ProjectFile.FileUrl = updateProjectFile.FileUrl
		response.ProjectFile.FileName = updateProjectFile.FileName
	}

	return response, nil
}

func (t *taskAndOwnerService) UpdateValidationOwner(taskID uint, userID uint) error {
	err := t.taskAndOwnerRepository.UpdateValidationOwner(taskID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (t *taskAndOwnerService) UpdateValidationManager(taskID uint, userID uint) error {
	err := t.taskAndOwnerRepository.UpdateValidationManager(taskID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (t *taskAndOwnerService) UpdateValidationEmployee(taskID uint, userID uint) error {
	err := t.taskAndOwnerRepository.UpdateValidationEmployee(taskID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (t *taskAndOwnerService) DeleteManager(taskId uint, managerId uint) error {
	db, countEmployee, countPlanningFile, countProjectFile, err := t.taskAndOwnerRepository.DeleteManager(taskId, managerId)
	if err != nil {
		return err
	}

	var manager domain.Manager
	err = helper.ResetAutoIncrement(db, &manager, "id", "managers")
	if err != nil {
		return err
	}

	if countEmployee == 1 {
		var employee domain.Employee
		err = helper.ResetAutoIncrement(db, &employee, "id", "employees")
		if err != nil {
			return err
		}
	}

	if countPlanningFile == 1 {
		var planningFile domain.PlanningFile
		err = helper.ResetAutoIncrement(db, &planningFile, "id", "planning_files")
		if err != nil {
			return err
		}
	}

	if countProjectFile == 1 {
		var projectFile domain.ProjectFile
		err = helper.ResetAutoIncrement(db, &projectFile, "id", "project_files")
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *taskAndOwnerService) DeleteEmployee(taskId uint, employeeId uint) error {
	db, countProjectFile, err := t.taskAndOwnerRepository.DeleteEmployee(taskId, employeeId)
	if err != nil {
		return err
	}

	var employee domain.Employee
	err = helper.ResetAutoIncrement(db, &employee, "id", "employees")
	if err != nil {
		return err
	}

	if countProjectFile > 0 {
		var projectFile domain.ProjectFile
		err = helper.ResetAutoIncrement(db, &projectFile, "id", "project_files")
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *taskAndOwnerService) DeletePlanningFile(fileId uint) (string, error) {
	db, fileName, err := t.taskAndOwnerRepository.DeletePlanningFile(fileId)
	if err != nil {
		return "", err
	}

	var planningFile domain.PlanningFile
	err = helper.ResetAutoIncrement(db, &planningFile, "id", "planning_files")
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (t *taskAndOwnerService) DeleteProjectFile(fileId uint) (string, error) {
	db, fileName, err := t.taskAndOwnerRepository.DeleteProjectFile(fileId)
	if err != nil {
		return "", err
	}

	var projectFile domain.ProjectFile
	err = helper.ResetAutoIncrement(db, &projectFile, "id", "project_files")
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (t *taskAndOwnerService) DeleteTaskAndOwner(taskID uint) error {
	db, countOwners, countManager, countEmployee, countPlanningFile, countProjectFile, err := t.taskAndOwnerRepository.Delete(taskID)
	if err != nil {
		return err
	}
	if err == nil {
		err = helper.SetupS3DeleteAll()
		if err != nil {
			return err
		}
	}

	// reset auto increment
	if countOwners > 0 {
		var owner domain.Owner
		err = helper.ResetAutoIncrement(db, &owner, "id", "owners")
		if err != nil {
			return err
		}
	}

	if countManager > 0 {
		var manager domain.Manager
		err = helper.ResetAutoIncrement(db, &manager, "id", "managers")
		if err != nil {
			return err
		}
	}

	if countEmployee > 0 {
		var employee domain.Employee
		err = helper.ResetAutoIncrement(db, &employee, "id", "employees")
		if err != nil {
			return err
		}
	}

	if countPlanningFile > 0 {
		var planningFile domain.PlanningFile
		err = helper.ResetAutoIncrement(db, &planningFile, "id", "planning_files")
		if err != nil {
			return err
		}
	}

	if countProjectFile > 0 {
		var projectFile domain.ProjectFile
		err = helper.ResetAutoIncrement(db, &projectFile, "id", "project_files")
		if err != nil {
			return err
		}
	}

	var task domain.Task
	err = helper.ResetAutoIncrement(db, &task, "id", "tasks")
	if err != nil {
		return err
	}

	return nil
}
