package service

import (
	"errors"
	"fmt"
	"log"
	"manajemen_tugas_master/helper"
	"manajemen_tugas_master/model/domain"
	"manajemen_tugas_master/model/web"
	"manajemen_tugas_master/repository"
	"time"

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

	taskDB, ownerDB, err := t.taskAndOwnerRepository.Create(user, task, board)
	if err != nil {
		return nil, nil, err
	}

	return taskDB, ownerDB, nil
}

func (t *taskAndOwnerService) GetTaskAndOwnerById(id uint) (*domain.TaskWithInvitation, error) {
	return t.taskAndOwnerRepository.FindById(id)
}

func (t *taskAndOwnerService) FindAllTasksAndOwners() ([]*domain.TaskWithInvitation, error) {
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

func (t *taskAndOwnerService) UpdateTaskAndOwner(task *domain.Task, newManagers []domain.Manager, newEmployees []domain.Employee, planningFile *domain.PlanningFile, projectFile *domain.ProjectFile, taskID uint, boardID uint) (*web.UpdateResponse, error) {
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

	updateTask, updatePlanningFile, updateProjectFile, err := t.taskAndOwnerRepository.Update(task, planningFile, projectFile)
	if err != nil {
		return nil, err
	}

	// Persiapan respons
	response := &web.UpdateResponse{}
	var emailsSent []string

	// notif email
	if task.NameTask != "" {
		response.NameTask = updateTask.NameTask
		ownerEmail, managerEmails, employeeEmails, nametask, _, err := t.taskAndOwnerRepository.GetNameEmailsDescription(uint64(taskID))
		if err != nil {
			return nil, err
		}

		to := []string{ownerEmail}
		to = append(to, managerEmails...)
		to = append(to, employeeEmails...)

		subject := "Task Name Updated"
		body := helper.GetEmailTemplate("Name task Update", nametask, "Name Updated", fmt.Sprintf("The name of the task has been updated to '%s'.", updateTask.NameTask))

		err = helper.SendEmail(to, subject, body)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
		} else {
			emailsSent = append(emailsSent, "Name task Update, Email sent successfully")
		}

		log.Println(ownerEmail, managerEmails, employeeEmails)
	}

	response.PlanningDescription = updateTask.PlanningDescription

	// notif email
	if updateTask.PlanningStatus == "Approved" || updateTask.PlanningStatus == "Not Approved" {
		response.PlanningStatus = updateTask.PlanningStatus
		ownerEmail, managerEmails, employeeEmails, nametask, _, err := t.taskAndOwnerRepository.GetNameEmailsDescription(uint64(taskID))
		if err != nil {
			return nil, err
		}

		to := []string{ownerEmail}
		to = append(to, managerEmails...)
		to = append(to, employeeEmails...)

		var subject, body string
		if updateTask.PlanningStatus == "Approved" {
			subject = "Task Planning Approved"
			body = helper.GetEmailTemplate("Planning Status Update", nametask, "Approved", "The planning for this task has been approved.")
		} else {
			subject = "Task Planning Not Approved"
			body = helper.GetEmailTemplate("Planning Status Update", nametask, "Not Approved", "The planning for this task has not been approved. Please review and make necessary adjustments.")
		}

		err = helper.SendEmail(to, subject, body)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			return nil, fmt.Errorf("failed to send email: %v", err)
		} else {
			emailsSent = append(emailsSent, "Task Planning status Update Email sent successfully")
		}

		log.Println(ownerEmail, managerEmails, employeeEmails)
	}

	// notif email
	if updateTask.ProjectStatus == "Done" || updateTask.ProjectStatus == "Undone" || updateTask.ProjectStatus == "Working" {
		response.ProjectStatus = updateTask.ProjectStatus
		ownerEmail, managerEmails, employeeEmails, nametask, _, err := t.taskAndOwnerRepository.GetNameEmailsDescription(uint64(taskID))
		if err != nil {
			return nil, err
		}

		to := []string{ownerEmail}
		to = append(to, managerEmails...)
		to = append(to, employeeEmails...)

		var subject, body string
		switch updateTask.ProjectStatus {
		case "Done":
			subject = "Task Project Done"
			body = helper.GetEmailTemplate("Project Status Update", nametask, "Done", "The project for this task has been completed.")
		case "Undone":
			subject = "Task Project Undone"
			body = helper.GetEmailTemplate("Project Status Update", nametask, "Undone", "The project for this task has not been completed. Please review and make necessary adjustments.")
		case "Working":
			subject = "Task Project in Progress"
			body = helper.GetEmailTemplate("Project Status Update", nametask, "Working", "Work has started on the project for this task.")
		}

		err = helper.SendEmail(to, subject, body)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			return nil, fmt.Errorf("failed to send email: %v", err)
		} else {
			emailsSent = append(emailsSent, "Task Project status update, Email sent successfully")
		}

		log.Println(ownerEmail, managerEmails, employeeEmails)
	}

	if updateTask.PlanningDueDate != "" {
		response.PlanningDueDate = updateTask.PlanningDueDate

		_, managerEmails, _, TaskDescription, nametask, err := t.taskAndOwnerRepository.GetNameEmailsDescription(uint64(taskID))
		if err != nil {
			return nil, err
		}

		if len(managerEmails) > 0 {
			//Parse date from string to time.Time
			parsedDueDate, err := time.Parse("02-01-2006", updateTask.PlanningDueDate)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse planning due date: %v", err)
			}

			// Create Google Calendar event
			senderEmail := "manajementugasapp@gmail.com"
			summary := fmt.Sprintf("Task: %s", nametask)
			description := TaskDescription

			// Use current time for startDateTime
			startDateTime := time.Now().Format(time.RFC3339)

			// Format endDateTime as expected by Google Calendar API
			endDateTime := parsedDueDate.Format(time.RFC3339)

			timeZone := "Asia/Jakarta" // Adjust to desired timezone
			attendees := managerEmails

			event, err := helper.CreateGoogleCalendarEvent(senderEmail, summary, description, startDateTime, endDateTime, timeZone, attendees)
			if err != nil {
				log.Printf("Failed to create Google Calendar event: %v", err)
			} else {
				log.Printf("Google Calendar event created: %s", event.HtmlLink)
			}

			// Send email invitation
			emailSubject := fmt.Sprintf("Calendar Invite: %s", summary)
			emailBody := helper.GetCalendarInviteTemplate(summary, description)
			err = helper.SendEmail(attendees, emailSubject, emailBody)
			if err != nil {
				log.Printf("Failed to send email: %v", err)
			} else {
				emailsSent = append(emailsSent, "Task Planning due date Update, Email sent successfully")
			}
		}
	}

	// calendar schedule
	if updateTask.ProjectDueDate != "" {
		response.ProjectDueDate = updateTask.ProjectDueDate

		_, _, employeeEmails, TaskDescription, nametask, err := t.taskAndOwnerRepository.GetNameEmailsDescription(uint64(taskID))
		if err != nil {
			return nil, err
		}

		if len(employeeEmails) > 0 {
			//Parse date from string to time.Time
			parsedDueDate, err := time.Parse("02-01-2006", updateTask.ProjectDueDate)
			if err != nil {
				return nil, fmt.Errorf("Failed to parse project due date: %v", err)
			}

			// Create Google Calendar event
			senderEmail := "manajementugasapp@gmail.com"
			summary := fmt.Sprintf("Task: %s", nametask)
			description := TaskDescription

			// Use current time for startDateTime
			startDateTime := time.Now().Format(time.RFC3339)

			// Format endDateTime as expected by Google Calendar API
			endDateTime := parsedDueDate.Format(time.RFC3339)

			timeZone := "Asia/Jakarta" // Adjust to desired timezone
			attendees := employeeEmails

			event, err := helper.CreateGoogleCalendarEvent(senderEmail, summary, description, startDateTime, endDateTime, timeZone, attendees)
			if err != nil {
				log.Printf("Failed to create Google Calendar event: %v", err)
			} else {
				log.Printf("Google Calendar event created: %s", event.HtmlLink)
			}

			// Send email invitation
			emailSubject := fmt.Sprintf("Calendar Invite: %s", summary)
			emailBody := helper.GetCalendarInviteTemplate(summary, description)
			err = helper.SendEmail(attendees, emailSubject, emailBody)
			if err != nil {
				log.Printf("Failed to send email: %v", err)
			} else {
				emailsSent = append(emailsSent, "Task Project due date Update, Email sent successfully")
			}
		}
	}

	response.Priority = updateTask.Priority

	// notif email
	if updateTask.ProjectComment != "" {
		response.ProjectComment = updateTask.ProjectComment
		ownerEmail, _, _, nametask, _, err := t.taskAndOwnerRepository.GetNameEmailsDescription(uint64(taskID))
		if err != nil {
			return nil, err
		}

		to := []string{ownerEmail}

		subject := "New Project Comment Added"
		body := helper.GetEmailTemplate("Project Comment Update", nametask, "New Comment", fmt.Sprintf("A new comment has been added to the project:\n\n'%s'", updateTask.ProjectComment))

		err = helper.SendEmail(to, subject, body)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			return nil, fmt.Errorf("failed to send email: %v", err)
		} else {
			emailsSent = append(emailsSent, "Task Project comment Update Email sent successfully")
		}

		log.Println(ownerEmail)
	}

	// Add new managers
	for _, manager := range newManagers {
		managerResponse, err := t.AddManager(taskID, manager.Email)
		if err != nil {
			return nil, err
		}
		response.Managers = append(response.Managers, *managerResponse)
	}

	// Add new employees
	for _, employee := range newEmployees {
		employeeResponse, err := t.AddEmployee(taskID, employee.Email)
		if err != nil {
			return nil, err
		}
		response.Employees = append(response.Employees, *employeeResponse)
	}

	// Fetch updated managers and employees
	updatedManagers, err := t.taskAndOwnerRepository.GetManagersByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	updatedEmployees, err := t.taskAndOwnerRepository.GetEmployeesByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	// Prepare response
	response = &web.UpdateResponse{
		// ... (other fields)
		Managers:  make([]web.ManagerResponse, len(updatedManagers)),
		Employees: make([]web.EmployeeResponse, len(updatedEmployees)),
	}

	for i, manager := range updatedManagers {
		response.Managers[i] = web.ManagerResponse{
			ID:               manager.ID,
			Email:            manager.Email,
			CustomRole:       manager.CustomRole,
			InvitationID:     manager.InvitationID,
			InvitationStatus: manager.InvitationStatus,
		}
	}

	for i, employee := range updatedEmployees {
		response.Employees[i] = web.EmployeeResponse{
			ID:               employee.ID,
			Email:            employee.Email,
			CustomRole:       employee.CustomRole,
			InvitationID:     employee.InvitationID,
			InvitationStatus: employee.InvitationStatus,
		}
	}

	// Populate planningFileResponse dengan data dari updatePlanningFile jika tidak kosong
	if updatePlanningFile.ID != 0 || updatePlanningFile.FileUrl != "" || updatePlanningFile.FileName != "" {
		response.PlanningFile.ID = updatePlanningFile.ID
		response.PlanningFile.FileUrl = updatePlanningFile.FileUrl
		response.PlanningFile.FileName = updatePlanningFile.FileName

		// notif email
		ownerEmail, managerEmails, employeeEmails, nametask, _, err := t.taskAndOwnerRepository.GetNameEmailsDescription(uint64(taskID))
		if err != nil {
			return nil, err
		}

		to := []string{ownerEmail}
		to = append(to, managerEmails...)
		to = append(to, employeeEmails...)

		subject := "Planning File Updated"
		body := helper.GetEmailTemplate("Planning File Update", nametask, "File Updated", fmt.Sprintf("A planning file has been updated:\nFile Name: %s\nFile URL: %s", updatePlanningFile.FileName, updatePlanningFile.FileUrl))

		err = helper.SendEmail(to, subject, body)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			return nil, fmt.Errorf("failed to send email: %v", err)
		} else {
			emailsSent = append(emailsSent, "Task Planning file Update Email sent successfully")
		}

		log.Println(ownerEmail, managerEmails, employeeEmails)
	}

	// Populate projectFileResponse dengan data dari updateProjectFile jika tidak kosong
	if updateProjectFile.ID != 0 || updateProjectFile.FileUrl != "" || updateProjectFile.FileName != "" {
		response.ProjectFile.ID = updateProjectFile.ID
		response.ProjectFile.FileUrl = updateProjectFile.FileUrl
		response.ProjectFile.FileName = updateProjectFile.FileName

		// notif email
		ownerEmail, managerEmails, employeeEmails, nametask, _, err := t.taskAndOwnerRepository.GetNameEmailsDescription(uint64(taskID))
		if err != nil {
			return nil, err
		}

		to := []string{ownerEmail}
		to = append(to, managerEmails...)
		to = append(to, employeeEmails...)

		subject := "Project File Updated"
		body := helper.GetEmailTemplate("Project File Update", nametask, "File Updated", fmt.Sprintf("A project file has been updated:\nFile Name: %s\nFile URL: %s", updateProjectFile.FileName, updateProjectFile.FileUrl))

		err = helper.SendEmail(to, subject, body)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			return nil, fmt.Errorf("failed to send email: %v", err)
		} else {
			emailsSent = append(emailsSent, "Task Project file Update Email sent successfully")
		}

		log.Println(ownerEmail, managerEmails, employeeEmails)
	}

	response.EmailsSent = emailsSent

	return response, nil
}

func (t *taskAndOwnerService) UpdateOwnerCustomRole(taskID uint, customRole string) (*domain.Owner, error) {
	owner, err := t.taskAndOwnerRepository.UpdateOwnerCustomRole(taskID, customRole)
	if err != nil {
		return nil, err
	}
	return owner, nil
}

func (t *taskAndOwnerService) AddManager(taskID uint, email string) (*web.ManagerResponse, error) {
	manager, err := t.taskAndOwnerRepository.AddManager(taskID, email)
	if err != nil {
		return nil, err
	}

	return &web.ManagerResponse{
		ID:               manager.ID,
		Email:            manager.Email,
		CustomRole:       manager.CustomRole,
		InvitationID:     manager.InvitationID,
		InvitationStatus: manager.InvitationStatus,
	}, nil
}

func (t *taskAndOwnerService) AddEmployee(taskID uint, email string) (*web.EmployeeResponse, error) {
	employee, err := t.taskAndOwnerRepository.AddEmployee(taskID, email)
	if err != nil {
		return nil, err
	}

	return &web.EmployeeResponse{
		ID:               employee.ID,
		Email:            employee.Email,
		CustomRole:       employee.CustomRole,
		InvitationID:     employee.InvitationID,
		InvitationStatus: employee.InvitationStatus,
	}, nil
}

func (t *taskAndOwnerService) UpdateManagerEmail(taskID uint, oldEmail, newEmail string) (*web.ManagerResponse, error) {
	manager, err := t.taskAndOwnerRepository.UpdateManagerEmail(taskID, oldEmail, newEmail)
	if err != nil {
		return nil, err
	}

	return &web.ManagerResponse{
		ID:               manager.ID,
		Email:            manager.Email,
		CustomRole:       manager.CustomRole,
		InvitationID:     manager.InvitationID,
		InvitationStatus: manager.InvitationStatus,
	}, nil
}

func (t *taskAndOwnerService) UpdateEmployeeEmail(taskID uint, oldEmail, newEmail string) (*web.EmployeeResponse, error) {
	employee, err := t.taskAndOwnerRepository.UpdateEmployeeEmail(taskID, oldEmail, newEmail)
	if err != nil {
		return nil, err
	}

	return &web.EmployeeResponse{
		ID:               employee.ID,
		Email:            employee.Email,
		CustomRole:       employee.CustomRole,
		InvitationID:     employee.InvitationID,
		InvitationStatus: employee.InvitationStatus,
	}, nil
}

func (t *taskAndOwnerService) UpdateManagerCustomRole(taskID uint, email, customRole string) error {
	return t.taskAndOwnerRepository.UpdateManagerCustomRole(taskID, email, customRole)
}

func (t *taskAndOwnerService) UpdateEmployeeCustomRole(taskID uint, email, customRole string) error {
	return t.taskAndOwnerRepository.UpdateEmployeeCustomRole(taskID, email, customRole)
}

func (t *taskAndOwnerService) RespondToInvitation(invitationID uint64, response string, role string) (*domain.Invitation, error) {
	if response != "accept" && response != "reject" {
		return nil, errors.New("Invalid response")
	}

	if role != "manager" && role != "employee" {
		return nil, errors.New("Invalid role")
	}

	var err error
	var updatedInvitation *domain.Invitation
	if role == "manager" {
		updatedInvitation, err = t.taskAndOwnerRepository.UpdateManagerInvitationStatus(invitationID, response)
	} else {
		updatedInvitation, err = t.taskAndOwnerRepository.UpdateEmployeeInvitationStatus(invitationID, response)
	}

	if err != nil {
		return nil, err
	}

	return updatedInvitation, nil
}

func (t *taskAndOwnerService) GetManagersByTaskID(taskID uint) ([]web.ManagerResponse, error) {
	managers, err := t.taskAndOwnerRepository.GetManagersByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	var managerResponses []web.ManagerResponse
	for _, manager := range managers {
		managerResponses = append(managerResponses, web.ManagerResponse{
			ID:               manager.ID,
			Email:            manager.Email,
			CustomRole:       manager.CustomRole,
			InvitationID:     manager.InvitationID,
			InvitationStatus: manager.InvitationStatus,
		})
	}

	return managerResponses, nil
}

func (t *taskAndOwnerService) GetEmployeesByTaskID(taskID uint) ([]web.EmployeeResponse, error) {
	employees, err := t.taskAndOwnerRepository.GetEmployeesByTaskID(taskID)
	if err != nil {
		return nil, err
	}

	var employeeResponses []web.EmployeeResponse
	for _, employee := range employees {
		employeeResponses = append(employeeResponses, web.EmployeeResponse{
			ID:               employee.ID,
			Email:            employee.Email,
			CustomRole:       employee.CustomRole,
			InvitationID:     employee.InvitationID,
			InvitationStatus: employee.InvitationStatus,
		})
	}

	return employeeResponses, nil
}

func (t *taskAndOwnerService) GetAllInvitations() ([]domain.Invitation, error) {
	return t.taskAndOwnerRepository.GetAllInvitations()
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
