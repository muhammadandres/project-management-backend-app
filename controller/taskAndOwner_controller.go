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
	userCtx := ctx.Locals("user")
	var user *domain.User
	// melakukan type assertion untuk mengkonversi user yang tipe datanya sudah menjadi interface{} ke *domain.user, agar dapat mengakses data user.
	user = userCtx.(*domain.User)

	// name_task
	var task domain.Task
	nameTask := ctx.FormValue("name_task")
	task.NameTask = nameTask

	taskDB, ownerDB, err := t.taskAndOwnerService.CreateTaskAndOwner(user, &task)
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

	response := web.WebResponse{
		Code:    200,
		Message: "Success",
		Data: CreateResponse{
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

	// user yang sedang login
	user := ctx.Locals("user")
	userID := user.(*domain.User).ID

	// task id
	taskId := ctx.Params("id")
	taskIdUint64, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task Id"})
	}

	// manager
	managerEmail := ctx.FormValue("manager")
	//if len(managerEmail) == 0 {
	//	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form request"})
	//}
	if managerEmail != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		manager.Email = managerEmail
	}

	// employee
	employeeEmail := ctx.FormValue("employee")
	//if len(employeeEmail) == 0 {
	//	return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form request"})
	//}
	if employeeEmail != "" {
		if err := t.taskAndOwnerService.UpdateValidationManager(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		employee.Email = employeeEmail
	}

	// planning file
	planningFiles, err := ctx.FormFile("planning_file")
	// Validasi dengan UpdateTaskAndOwnerValidationForManager jika planning file diunggah terlebih dahulu
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
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error uploading project file" + err.Error()})
		}
	}

	// project file
	projectFiles, err := ctx.FormFile("project_file")
	// Validasi dengan UpdateTaskAndOwnerValidationForEmployee jika project file diunggah terlebih dahulu
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
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error uploading project file" + err.Error()})
		}
	}

	// field yang tidak berelasi pada task
	nameTask := ctx.FormValue("name_task")
	if nameTask != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.NameTask = nameTask
	}

	planningDescription := ctx.FormValue("planning_description")
	if planningDescription != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.PlanningDescription = planningDescription
	}

	planningStatus := ctx.FormValue("planning_status")
	if planningStatus != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.PlanningStatus = planningStatus

	}

	projectStatus := ctx.FormValue("project_status")
	if projectStatus != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.ProjectStatus = projectStatus
	}

	planningDueDate := ctx.FormValue("planning_due_date")
	if planningDueDate != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.PlanningDueDate = planningDueDate
	}

	projectDueDate := ctx.FormValue("project_due_date")
	if projectDueDate != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.ProjectDueDate = projectDueDate
	}

	priority := ctx.FormValue("priority")
	if priority != "" {
		if err := t.taskAndOwnerService.UpdateValidationManager(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.Priority = priority
	}

	projectComment := ctx.FormValue("project_comment")
	if projectComment != "" {
		if err := t.taskAndOwnerService.UpdateValidationEmployee(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		task.ProjectComment = projectComment
	}

	// save
	response, err := t.taskAndOwnerService.UpdateTaskAndOwner(&task, &manager, &employee, &planningFile, &projectFile, uint(taskIdUint64))
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
	user := ctx.Locals("user")
	userId := user.(*domain.User).ID

	taskId := ctx.Params("id")
	taskIdUint64, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid file Id"})
	}

	managerId := ctx.Params("manager_id")
	managerIdUint64, err := strconv.ParseUint(managerId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task Id"})
	}

	if userId != 0 && taskId != "" {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userId)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "task dan manager id is required"})
	}

	if err := t.taskAndOwnerService.DeleteManager(uint(taskIdUint64), uint(managerIdUint64)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "File deleted successfully"})
}

func (t *TaskAndOwnerController) DeleteEmployee(ctx *fiber.Ctx) error {
	user := ctx.Locals("user")
	userId := user.(*domain.User).ID

	taskId := ctx.Params("id")
	taskIdUint64, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task Id"})
	}

	employeeId := ctx.Params("employee_id")
	employeeIdUint64, err := strconv.ParseUint(employeeId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid file Id"})
	}

	if userId != 0 && taskId != "" {
		if err := t.taskAndOwnerService.UpdateValidationManager(uint(taskIdUint64), uint(userId)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "task and employee id is required"})
	}

	if err := t.taskAndOwnerService.DeleteEmployee(uint(taskIdUint64), uint(employeeIdUint64)); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "File deleted successfully"})
}

func (t *TaskAndOwnerController) DeletePlanningFile(ctx *fiber.Ctx) error {
	user := ctx.Locals("user")
	userId := user.(*domain.User).ID

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

	if userId != 0 && taskId != "" {
		if err := t.taskAndOwnerService.UpdateValidationManager(uint(taskIdUint64), uint(userId)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "task and file id is required"})
	}

	fileName, err := t.taskAndOwnerService.DeletePlanningFile(uint(fileIdUint64))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// aws s3
	err = helper.SetupS3Delete(fileName)
	if err != nil {
		// Periksa tipe error untuk penanganan khusus
		switch err.(type) {
		case error: // pesan error yang bertipe error di ambil dari fmt.Errorf
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		default:
			// pesan error yang bukan bertipe error di ambil dari errors.New
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "File deleted successfully"})
}

func (t *TaskAndOwnerController) DeleteProjectFile(ctx *fiber.Ctx) error {
	user := ctx.Locals("user")
	userId := user.(*domain.User).ID

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

	if userId != 0 && taskId != "" {
		if err := t.taskAndOwnerService.UpdateValidationEmployee(uint(taskIdUint64), uint(userId)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "task and file id is required"})
	}

	fileName, err := t.taskAndOwnerService.DeleteProjectFile(uint(fileIdUint64))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// aws s3
	err = helper.SetupS3Delete(fileName)
	if err != nil {
		// Periksa tipe error untuk penanganan khusus
		switch err.(type) {
		case error: // pesan error yang bertipe error di ambil dari fmt.Errorf
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		default:
			// pesan error yang bukan bertipe error di ambil dari errors.New
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "File deleted successfully"})
}

func (t *TaskAndOwnerController) DeleteTaskAndOwner(ctx *fiber.Ctx) error {
	taskId := ctx.Params("id")
	taskIdUint64, err := strconv.ParseUint(taskId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid task Id"})
	}

	user := ctx.Locals("user")
	userID := user.(*domain.User).ID

	if taskId != "" && userID != 0 {
		if err := t.taskAndOwnerService.UpdateValidationOwner(uint(taskIdUint64), uint(userID)); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
	} else {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "task and file id is required"})
	}

	if err := t.taskAndOwnerService.DeleteTaskAndOwner(uint(taskIdUint64)); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Deleted successfully"})
}
