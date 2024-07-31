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

	// board initialize
	boardRepository, _ := InitializeRepositoryBoard(db)
	boardService, _ := InitializeServiceBoard(boardRepository)
	boardController, _ := InitializeControllerBoard(boardService)

	// task initialize
	taskRepository, _ := InitializeRepositoryTask(db)
	taskService, _ := InitializeServiceTask(taskRepository, boardRepository)
	taskController, _ := InitializeControllerTask(taskService)

	// app.Post("/test/calendar", func(c *fiber.Ctx) error {
	// 	// Data contoh
	// 	senderEmail := "m.andres.novrizal@gmail.com"
	// 	summary := "Test Event"
	// 	description := "This is a test event created via API"
	// 	startDateTime := time.Now().Format(time.RFC3339)
	// 	endDateTime := time.Now().Add(1 * time.Hour).Format(time.RFC3339)
	// 	timeZone := "Asia/Jakarta"
	// 	attendees := []string{"m.andres.novrizal@gmail.com"}

	// 	// Membuat acara kalender
	// 	event, err := helper.CreateGoogleCalendarEvent(senderEmail, summary, description, startDateTime, endDateTime, timeZone, attendees)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": "Failed to create calendar event: " + err.Error(),
	// 		})
	// 	}

	// 	emailSubject := "Test Calendar Invite"
	// 	emailBody := "You've been invited to an event. Check your calendar."
	// 	err = helper.SendEmail(attendees, emailSubject, emailBody)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 			"error": "Failed to send email: " + err.Error(),
	// 		})
	// 	}

	// 	return c.JSON(fiber.Map{
	// 		"message":   "Calendar event created and email sent",
	// 		"eventLink": event.HtmlLink,
	// 	})
	// })

	// Group route untuk user
	userRoutes := app.Group("/")
	userRoutes.Post("user/signup", userController.SignupUser)
	userRoutes.Post("user/login", userController.LoginUser)
	userRoutes.Get("users", userController.GetAllUsers)
	userRoutes.Get("auth/oauth", userController.GoogleOauth)
	userRoutes.Get("auth/callback", userController.GoogleCallback)
	userRoutes.Post("user/forgot-password", userController.ForgotPassword)
	userRoutes.Get("user/:id", userController.GetUserByID)
	userRoutes.Use(middleware.AuthUser(userService))
	userRoutes.Delete("user/:id", userController.DeleteUser)
	userRoutes.Put("user/:id", userController.UpdateUser)

	// Groupt route untuk board
	boardRoutes := app.Group("/")
	boardRoutes.Use(middleware.AuthUser(userService)) // Gunakan middleware untuk semua route dalam grup task
	boardRoutes.Post("board", boardController.CreateBoard)
	boardRoutes.Get("board/:id", boardController.GetBoardById)
	boardRoutes.Get("boards", boardController.GetAllBoards)
	boardRoutes.Put("board/:id", boardController.EditBoard)
	boardRoutes.Delete("board/:id", boardController.DeleteBoardById)

	// Group route untuk task
	taskRoutes := app.Group("/")
	taskRoutes.Use(middleware.AuthUser(userService)) // Gunakan middleware untuk semua route dalam grup task
	taskRoutes.Post("task/:board_id", taskController.CreateTaskAndOwner)
	taskRoutes.Put("board/:boardId/task/:taskId", taskController.UpdateTaskAndOwner)
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
