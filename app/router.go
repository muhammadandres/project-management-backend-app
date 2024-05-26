package app

import (
	"manajemen_tugas_master/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// user initialize
	userRepository, _ := InitializeRepositoryUser(db)
	userService, _ := InitializeServiceUser(userRepository)
	userController, _ := InitializeControllerUser(userService)

	// task initialize
	taskRepository, _ := InitializeRepositoryTask(db)
	taskService, _ := InitializeServiceTask(taskRepository)
	taskController, _ := InitializeControllerTask(taskService)

	// Group route untuk user
	userRoutes := app.Group("/")
	userRoutes.Post("user/signup", userController.SignupUser)
	userRoutes.Post("user/login", userController.LoginUser)
	userRoutes.Get("users", userController.GetAllUsers)
	userRoutes.Get("user/:id", userController.GetUserByID)
	userRoutes.Delete("user/:id", userController.DeleteUser)
	userRoutes.Use(middleware.RequireAuthUser(userService)) // Gunakan middleware untuk semua route dalam grup user
	userRoutes.Put("user/:id", userController.UpdateUser)

	// Group route untuk task
	taskRoutes := app.Group("/")
	taskRoutes.Use(middleware.RequireAuthUser(userService)) // Gunakan middleware untuk semua route dalam grup task
	taskRoutes.Post("task", taskController.CreateTaskAndOwner)
	taskRoutes.Put("task/:id", taskController.UpdateTaskAndOwner)
	taskRoutes.Get("task/:id", taskController.GetTaskAndOwnerById)
	taskRoutes.Get("tasks", taskController.GetAllTasksAndOwners)
	taskRoutes.Get("tasks/owners", taskController.GetAllOwners)
	taskRoutes.Get("tasks/managers", taskController.GetAllManagers)
	taskRoutes.Get("tasks/employees", taskController.GetAllEmployees)
	taskRoutes.Get("tasks/planning_files", taskController.GetAllPlanningFiles)
	taskRoutes.Get("tasks/project_files", taskController.GetAllProjectFiles)
	taskRoutes.Delete("task/:id/manager/:manager_id", taskController.DeleteManager)
	taskRoutes.Delete("task/:id/employee/:employee_id", taskController.DeleteEmployee)
	taskRoutes.Delete("task/:id/planning_file/:file_id", taskController.DeletePlanningFile)
	taskRoutes.Delete("task/:id/project_file/:file_id", taskController.DeleteProjectFile)
	taskRoutes.Delete("task/:id", taskController.DeleteTaskAndOwner)
}
