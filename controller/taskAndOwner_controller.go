package controller

import (
	"manajemen_tugas_master/helper"
	"manajemen_tugas_master/model/domain"
	"manajemen_tugas_master/model/web"
	"manajemen_tugas_master/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TaskAndOwnerController struct {
	taskAndOwnerService service.TaskAndOwnerService
}

func NewTaskController(taskAndOwnerService service.TaskAndOwnerService) *TaskAndOwnerController {
	return &TaskAndOwnerController{taskAndOwnerService}
}

func (t *TaskAndOwnerController) CreateTaskAndOwner(ctx *fiber.Ctx) error {
	var user *domain.User
	userCtx := ctx.Locals("user")
	if userCtx != nil {
		var ok bool
		user, ok = userCtx.(*domain.User)
		if !ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user data"})
		}
	}

	userOauth := ctx.Locals("userOauth")
	if userOauth != nil {
		var ok bool
		user, ok = userOauth.(*domain.User)
		if !ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid OAuth user data"})
		}
	}

	// Dapatkan board ID dari parameter
	boardIDStr := ctx.Params("board_id")
	if boardIDStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID board diperlukan"})
	}

	boardID, err := strconv.ParseUint(boardIDStr, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID board tidak valid"})
	}

	// board_id
	board := &domain.Board{ID: boardID}

	// name_task
	nameTask := ctx.FormValue("name_task")
	task := &domain.Task{NameTask: nameTask}

	taskDB, ownerDB, err := t.taskAndOwnerService.CreateTaskAndOwner(user, task, board)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	type CreateResponse struct {
		BoardId   uint64 `json:"board_id"`
		TaskID    uint64 `json:"task_id"`
		NameTask  string `json:"name_task"`
		OwnerID   uint64 `json:"owner_id"`
		UserEmail string `json:"user_email"`
		UserID    uint64 `json:"user_id"`
	}

	response := web.WebResponse{
		Code:    200,
		Message: "Success",
		Data: CreateResponse{
			BoardId:   taskDB.BoardID,
			TaskID:    taskDB.ID,
			NameTask:  taskDB.NameTask,
			OwnerID:   ownerDB.ID,
			UserEmail: ownerDB.Email,
			UserID:    ownerDB.UserID,
		},
	}

	return ctx.Status(fiber.StatusCreated).JSON(response)
}

func (t *TaskAndOwnerController) GetTaskAndOwnerById(ctx *fiber.Ctx) error {
	taskId := ctx.Params("id")
	taskIdUint64, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task Id"})
	}

	task, err := t.taskAndOwnerService.GetTaskAndOwnerById(uint(taskIdUint64))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Task not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(web.CreateResponseTask(task))
}

func (t *TaskAndOwnerController) GetAllTasksAndOwners(ctx *fiber.Ctx) error {
	tasks, err := t.taskAndOwnerService.FindAllTasksAndOwners()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(web.CreateResponseTasks(tasks))
}

func (t *TaskAndOwnerController) GetAllOwners(ctx *fiber.Ctx) error {
	tasks, err := t.taskAndOwnerService.FindAllOwners()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	type CreateResponse struct {
		TaskID    uint64 `json:"task_id"`
		NameTask  string `json:"name_task"`
		OwnerID   uint64 `json:"owner_id"`
		UserEmail string `json:"user_email"`
		UserID    uint64 `json:"user_id"`
	}

	var response []web.WebResponse
	for _, task := range tasks {
		response = append(response, web.WebResponse{
			Code:    200,
			Message: "Success",
			Data: CreateResponse{
				TaskID:    task.ID,
				NameTask:  task.NameTask,
				OwnerID:   task.Owner.ID,
				UserEmail: task.Owner.Email,
				UserID:    task.Owner.UserID,
			},
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (t *TaskAndOwnerController) GetAllManagers(ctx *fiber.Ctx) error {
	tasks, err := t.taskAndOwnerService.FindAllManagers()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	type CreateResponse struct {
		TaskID    uint64 `json:"task_id"`
		NameTask  string `json:"name_task"`
		ManagerID uint64 `json:"owner_id"`
		UserEmail string `json:"user_email"`
		UserID    uint64 `json:"user_id"`
	}

	var response []web.WebResponse
	for _, task := range tasks {
		for _, manager := range task.Manager {
			response = append(response, web.WebResponse{
				Code:    200,
				Message: "Success",
				Data: CreateResponse{
					TaskID:    task.ID,
					NameTask:  task.NameTask,
					ManagerID: manager.ID,
					UserEmail: manager.Email,
					UserID:    manager.UserID,
				},
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (t *TaskAndOwnerController) GetAllEmployees(ctx *fiber.Ctx) error {
	tasks, err := t.taskAndOwnerService.FindAllEmployees()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	type CreateResponse struct {
		TaskID     uint64 `json:"task_id"`
		NameTask   string `json:"name_task"`
		EmployeeID uint64 `json:"owner_id"`
		UserEmail  string `json:"user_email"`
		UserID     uint64 `json:"user_id"`
	}

	var response []web.WebResponse
	for _, task := range tasks {
		for _, employee := range task.Employee {
			response = append(response, web.WebResponse{
				Code:    200,
				Message: "Success",
				Data: CreateResponse{
					TaskID:     task.ID,
					NameTask:   task.NameTask,
					EmployeeID: employee.ID,
					UserEmail:  employee.Email,
					UserID:     employee.UserID,
				},
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (t *TaskAndOwnerController) GetAllPlanningFiles(ctx *fiber.Ctx) error {
	tasks, err := t.taskAndOwnerService.FindAllPlanningFiles()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	type CreateResponse struct {
		TaskID         uint64 `json:"task_id"`
		NameTask       string `json:"name_task"`
		PlanningFileID uint64 `json:"planning_file_id"`
		FileUrl        string `json:"file_url"`
		FileName       string `json:"file_name"`
	}

	var response []web.WebResponse
	for _, task := range tasks {
		for _, file := range task.PlanningFile {
			response = append(response, web.WebResponse{
				Code:    200,
				Message: "Success",
				Data: CreateResponse{
					TaskID:         task.ID,
					NameTask:       task.NameTask,
					PlanningFileID: file.ID,
					FileUrl:        file.FileUrl,
					FileName:       file.FileName,
				},
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (t *TaskAndOwnerController) GetAllProjectFiles(ctx *fiber.Ctx) error {
	tasks, err := t.taskAndOwnerService.FindAllProjectFiles()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	type CreateResponse struct {
		TaskID        uint64 `json:"task_id"`
		NameTask      string `json:"name_task"`
		ProjectFileID uint64 `json:"project_file_id"`
		FileUrl       string `json:"file_url"`
		FileName      string `json:"file_name"`
	}

	var response []web.WebResponse
	for _, task := range tasks {
		for _, file := range task.ProjectFile {
			response = append(response, web.WebResponse{
				Code:    200,
				Message: "Success",
				Data: CreateResponse{
					TaskID:        task.ID,
					NameTask:      task.NameTask,
					ProjectFileID: file.ID,
					FileUrl:       file.FileUrl,
					FileName:      file.FileName,
				},
			})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (t *TaskAndOwnerController) UpdateTaskAndOwner(ctx *fiber.Ctx) error {
	var (
		task         domain.Task
		planningFile domain.PlanningFile
		projectFile  domain.ProjectFile
		manager      domain.Manager
		employee     domain.Employee
	)

	// Get user from context (either JWT or OAuth)
	var userID uint64
	user := ctx.Locals("user")
	userOauth := ctx.Locals("userOauth")

	if user != nil {
		userID = user.(*domain.User).ID
	} else if userOauth != nil {
		userID = userOauth.(*domain.User).ID
	}

	// board id
	boardId := ctx.Params("boardId")
	boardIdUint64, err := strconv.ParseUint(boardId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid board Id"})
	}

	// task id
	taskId := ctx.Params("taskId")
	taskIdUint64, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// manager
	managerEmail := ctx.FormValue("manager")
	if managerEmail != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		manager.Email = managerEmail
	}

	// employee
	employeeEmail := ctx.FormValue("employee")
	if employeeEmail != "" {
		if err := t.taskAndOwnerService.UpdateValidationManager(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		employee.Email = employeeEmail
	}

	// planning file
	planningFiles, err := ctx.FormFile("planning_file")
	if planningFiles != nil {
		if err := t.taskAndOwnerService.UpdateValidationManager(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err == nil {
			PlanningFileUrl, PlanningFileName, err := helper.SetupS3Uploader(planningFiles)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error uploading planning file" + err.Error()})
			}
			planningFile.FileUrl = PlanningFileUrl
			planningFile.FileName = PlanningFileName
		}
	}

	// project file
	projectFiles, err := ctx.FormFile("project_file")
	if projectFiles != nil {
		if err := t.taskAndOwnerService.UpdateValidationEmployee(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err == nil {
			ProjectFileUrl, ProjectFileName, err := helper.SetupS3Uploader(projectFiles)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error uploading project file" + err.Error()})
			}
			projectFile.FileUrl = ProjectFileUrl
			projectFile.FileName = ProjectFileName
		}
	}

	// field yang tidak berelasi pada task
	if nameTask := ctx.FormValue("name_task"); nameTask != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.NameTask = nameTask
	}

	if planningDescription := ctx.FormValue("planning_description"); planningDescription != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.PlanningDescription = planningDescription
	}

	if planningStatus := ctx.FormValue("planning_status"); planningStatus != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.PlanningStatus = planningStatus
	}

	if projectStatus := ctx.FormValue("project_status"); projectStatus != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.ProjectStatus = projectStatus
	}

	if planningDueDate := ctx.FormValue("planning_due_date"); planningDueDate != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.PlanningDueDate = planningDueDate
	}

	if projectDueDate := ctx.FormValue("project_due_date"); projectDueDate != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.ProjectDueDate = projectDueDate
	}

	if priority := ctx.FormValue("priority"); priority != "" {
		if err := t.taskAndOwnerService.UpdateValidationManager(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.Priority = priority
	}

	if projectComment := ctx.FormValue("project_comment"); projectComment != "" {
		if err := t.taskAndOwnerService.UpdateValidationEmployee(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.ProjectComment = projectComment
	}

	// save
	response, err := t.taskAndOwnerService.UpdateTaskAndOwner(&task, &manager, &employee, &planningFile, &projectFile, uint(taskIdUint64), uint(boardIdUint64))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    200,
		Message: "Success",
		Data:    response,
	})
}

func (t *TaskAndOwnerController) DeleteManager(ctx *fiber.Ctx) error {
	var userId, userOauthId uint64

	if user := ctx.Locals("user"); user != nil {
		if u, ok := user.(*domain.User); ok {
			userId = u.ID
		}
	}

	if userOauth := ctx.Locals("userOauth"); userOauth != nil {
		if u, ok := userOauth.(*domain.User); ok {
			userOauthId = u.ID
		}
	}

	taskId := ctx.Params("id")
	taskIdUint64, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task Id"})
	}

	managerId := ctx.Params("manager_id")
	managerIdUint64, err := strconv.ParseUint(managerId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid manager Id"})
	}

	if userId != 0 && taskId != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userId)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	} else if userOauthId != 0 && taskId != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userOauthId)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "task and manager id are required"})
	}

	if err := t.taskAndOwnerService.DeleteManager(uint(taskIdUint64), uint(managerIdUint64)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Manager deleted successfully"})
}

func (t *TaskAndOwnerController) DeleteEmployee(ctx *fiber.Ctx) error {
	taskId := ctx.Params("id")
	taskIdUint64, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task Id"})
	}

	employeeId := ctx.Params("employee_id")
	employeeIdUint64, err := strconv.ParseUint(employeeId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid employee Id"})
	}

	// Check for user authentication
	var authenticatedUserId uint
	if user := ctx.Locals("user"); user != nil {
		if u, ok := user.(*domain.User); ok {
			authenticatedUserId = uint(u.ID)
		}
	}

	if userOauth := ctx.Locals("userOauth"); userOauth != nil {
		if u, ok := userOauth.(*domain.User); ok {
			authenticatedUserId = uint(u.ID)
		}
	}

	if authenticatedUserId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	// Validate manager
	if err := t.taskAndOwnerService.UpdateValidationManager(uint(taskIdUint64), authenticatedUserId); err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	// Delete employee
	if err := t.taskAndOwnerService.DeleteEmployee(uint(taskIdUint64), uint(employeeIdUint64)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Employee deleted successfully"})
}

func (t *TaskAndOwnerController) DeletePlanningFile(ctx *fiber.Ctx) error {
	var authenticatedUserId uint64

	if user := ctx.Locals("user"); user != nil {
		if u, ok := user.(*domain.User); ok {
			authenticatedUserId = u.ID
		}
	} else if userOauth := ctx.Locals("userOauth"); userOauth != nil {
		if u, ok := userOauth.(*domain.User); ok {
			authenticatedUserId = u.ID
		}
	}

	if authenticatedUserId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	taskId := ctx.Params("id")
	taskIdUint64, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task Id"})
	}

	fileId := ctx.Params("file_id")
	fileIdUint64, err := strconv.ParseUint(fileId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid file Id"})
	}

	// Validate manager
	if err := t.taskAndOwnerService.UpdateValidationManager(uint(taskIdUint64), uint(authenticatedUserId)); err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	fileName, err := t.taskAndOwnerService.DeletePlanningFile(uint(fileIdUint64))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// AWS S3 delete
	err = helper.SetupS3Delete(fileName)
	if err != nil {
		if _, ok := err.(error); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "File deleted successfully"})
}

func (t *TaskAndOwnerController) DeleteProjectFile(ctx *fiber.Ctx) error {
	var authenticatedUserId uint64

	if user := ctx.Locals("user"); user != nil {
		if u, ok := user.(*domain.User); ok {
			authenticatedUserId = u.ID
		}
	} else if userOauth := ctx.Locals("userOauth"); userOauth != nil {
		if u, ok := userOauth.(*domain.User); ok {
			authenticatedUserId = u.ID
		}
	}

	if authenticatedUserId == 0 {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	taskId := ctx.Params("id")
	taskIdUint64, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task Id"})
	}

	fileId := ctx.Params("file_id")
	fileIdUint64, err := strconv.ParseUint(fileId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid file Id"})
	}

	// Validate employee
	if err := t.taskAndOwnerService.UpdateValidationEmployee(uint(taskIdUint64), uint(authenticatedUserId)); err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
	}

	fileName, err := t.taskAndOwnerService.DeleteProjectFile(uint(fileIdUint64))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// AWS S3 delete
	err = helper.SetupS3Delete(fileName)
	if err != nil {
		if _, ok := err.(error); ok {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "File deleted successfully"})
}

func (t *TaskAndOwnerController) DeleteTaskAndOwner(ctx *fiber.Ctx) error {
	taskId := ctx.Params("id")
	taskIdUint64, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task Id"})
	}

	var userId, userOauthId uint64

	if user := ctx.Locals("user"); user != nil {
		if u, ok := user.(*domain.User); ok {
			userId = u.ID
		}
	}

	if userOauth := ctx.Locals("userOauth"); userOauth != nil {
		if u, ok := userOauth.(*domain.User); ok {
			userOauthId = u.ID
		}
	}

	if taskId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "task id is required"})
	}

	if userId != 0 {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userId)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	} else if userOauthId != 0 {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userOauthId)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User authentication is required"})
	}

	if err := t.taskAndOwnerService.DeleteTaskAndOwner(uint(taskIdUint64)); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Deleted successfully"})
}
