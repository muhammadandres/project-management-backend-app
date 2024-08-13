package app

import (
	"manajemen_tugas_master/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB, store *session.Store) {
	// user initialize
	userRepository, _ := InitializeRepositoryUser(db)
	userService, _ := InitializeServiceUser(userRepository)
	userController, _ := InitializeControllerUser(userService, store)

	// board initialize
	boardRepository, _ := InitializeRepositoryBoard(db)
	boardService, _ := InitializeServiceBoard(boardRepository)
	boardController, _ := InitializeControllerBoard(boardService)

	// task initialize
	taskRepository, _ := InitializeRepositoryTask(db)
	taskService, _ := InitializeServiceTask(taskRepository, boardRepository)
	taskController, _ := InitializeControllerTask(taskService)

	app.Get("/", func(c *fiber.Ctx) error {
		tokenStringJwt := c.Cookies("Authorization")
		tokenStringOauth := c.Cookies("GoogleAuthorization")

		// Check if either JWT or OAuth token is present
		if tokenStringJwt != "" || tokenStringOauth != "" {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "You are authorized"})
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "You are not authorized. Please sign in or sign up to access this resource.",
		})
	})

	// Group route untuk user
	userRoutes := app.Group("/")
	userRoutes.Post("user/signup", userController.SignupUser)
	userRoutes.Post("user/login", userController.LoginUser)
	userRoutes.Get("users", userController.GetAllUsers)
	userRoutes.Get("auth/oauth", userController.GoogleOauth)
	userRoutes.Get("auth/callback", userController.GoogleCallback)
	userRoutes.Post("forgot-password", userController.ForgotPassword)
	userRoutes.Post("reset-password", userController.ResetPassword)
	userRoutes.Get("user/:id", userController.GetUserByID)
	userRoutes.Use(middleware.AuthUser(userService, store))
	userRoutes.Delete("user/:id", userController.DeleteUser)
	userRoutes.Put("user/:id", userController.UpdateUser)

	// Groupt route untuk board
	boardRoutes := app.Group("/")
	boardRoutes.Use(middleware.AuthUser(userService, store)) // Gunakan middleware untuk semua route dalam grup task
	boardRoutes.Post("board", boardController.CreateBoard)
	boardRoutes.Get("board/:id", boardController.GetBoardById)
	boardRoutes.Get("boards", boardController.GetAllBoards)
	boardRoutes.Put("board/:id", boardController.EditBoard)
	boardRoutes.Delete("board/:id", boardController.DeleteBoardById)

	// Group route untuk task
	taskRoutes := app.Group("/")
	taskRoutes.Use(middleware.AuthUser(userService, store)) // Gunakan middleware untuk semua route dalam grup task
	taskRoutes.Post("task/:board_id", taskController.CreateTaskAndOwner)
	taskRoutes.Put("board/:boardId/task/:taskId", taskController.UpdateTaskAndOwner)
	taskRoutes.Put("invitations/:invitationID/respond", taskController.RespondToInvitation)
	taskRoutes.Get("invitations", taskController.GetAllInvitations)
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
