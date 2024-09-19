package controller

import (
	"manajemen_tugas_master/model/domain"
	"manajemen_tugas_master/service"

	"github.com/gofiber/fiber/v2"
)

type BoardController struct {
	boardService service.BoardService
}

func NewBoardController(boardService service.BoardService) *BoardController {
	return &BoardController{boardService}
}

// CreateBoard godoc
// @Summary Create a new board
// @Description Create a new board. This endpoint requires cookie authentication.
// @Tags boards
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param name_board formData string true "Name of the board"
// @Success 201 {object} web.BoardResponse
// @Failure 400 {object} web.ErrorResponse
// @Failure 401 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /board [post]
func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
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

	if user == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not authenticated"})
	}

	board := &domain.Board{}
	nameBoard := ctx.FormValue("name_board")
	board.NameBoard = nameBoard
	board.UserID = user.ID

	boardDb, err := c.boardService.CreateBoard(board)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":    "Board berhasil dibuat",
		"board_id":   boardDb.ID,
		"name_board": boardDb.NameBoard,
		"user_id":    boardDb.UserID,
		"user_email": user.Email,
	})
}

func (c *BoardController) GetBoardById(ctx *fiber.Ctx) error {
	boardID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid board ID"})
	}

	board, err := c.boardService.GetBoardById(uint64(boardID))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	// Create a custom response structure
	response := struct {
		ID        uint64 `json:"id"`
		NameBoard string `json:"name_board"`
		Tasks     []struct {
			ID                        uint64                           `json:"id"`
			BoardID                   uint64                           `json:"board_id"`
			Owner                     domain.Owner                     `json:"owner"`
			Manager                   []domain.ManagerWithInvitation   `json:"manager"`
			Employee                  []domain.EmployeeWithInvitation  `json:"employee"`
			NameTask                  string                           `json:"name_task"`
			PlanningDescriptionPersen string                           `json:"planning_description_persen"`
			PlanningDescriptionFile   []domain.PlanningDescriptionFile `json:"planning_description_files"`
			PlanningFile              []domain.PlanningFile            `json:"planning_file"`
			PlanningStatus            string                           `json:"planning_status"`
			ProjectFile               []domain.ProjectFile             `json:"project_file"`
			ProjectStatus             string                           `json:"project_status"`
			PlanningDueDate           string                           `json:"planning_due_date"`
			ProjectDueDate            string                           `json:"project_due_date"`
			Priority                  string                           `json:"priority"`
			ProjectComment            string                           `json:"project_comment"`
		} `json:"tasks"`
	}{
		ID:        board.ID,
		NameBoard: board.NameBoard,
	}

	for _, task := range board.Tasks {
		taskResponse := struct {
			ID                        uint64                           `json:"id"`
			BoardID                   uint64                           `json:"board_id"`
			Owner                     domain.Owner                     `json:"owner"`
			Manager                   []domain.ManagerWithInvitation   `json:"manager"`
			Employee                  []domain.EmployeeWithInvitation  `json:"employee"`
			NameTask                  string                           `json:"name_task"`
			PlanningDescriptionPersen string                           `json:"planning_description_persen"`
			PlanningDescriptionFile   []domain.PlanningDescriptionFile `json:"planning_description_files"`
			PlanningFile              []domain.PlanningFile            `json:"planning_file"`
			PlanningStatus            string                           `json:"planning_status"`
			ProjectFile               []domain.ProjectFile             `json:"project_file"`
			ProjectStatus             string                           `json:"project_status"`
			PlanningDueDate           string                           `json:"planning_due_date"`
			ProjectDueDate            string                           `json:"project_due_date"`
			Priority                  string                           `json:"priority"`
			ProjectComment            string                           `json:"project_comment"`
		}{
			ID:                        task.ID,
			BoardID:                   task.BoardID,
			Owner:                     task.Owner,
			NameTask:                  task.NameTask,
			PlanningDescriptionPersen: task.PlanningDescriptionPersen,
			PlanningDescriptionFile:   task.PlanningDescriptionFile,
			PlanningFile:              task.PlanningFile,
			PlanningStatus:            task.PlanningStatus,
			ProjectFile:               task.ProjectFile,
			ProjectStatus:             task.ProjectStatus,
			PlanningDueDate:           task.PlanningDueDate,
			ProjectDueDate:            task.ProjectDueDate,
			Priority:                  task.Priority,
			ProjectComment:            task.ProjectComment,
		}

		for _, manager := range task.Manager {
			taskResponse.Manager = append(taskResponse.Manager, domain.ManagerWithInvitation{
				Manager:          manager,
				InvitationStatus: manager.InvitationStatus,
			})
		}

		for _, employee := range task.Employee {
			taskResponse.Employee = append(taskResponse.Employee, domain.EmployeeWithInvitation{
				Employee:         employee,
				InvitationStatus: employee.InvitationStatus,
			})
		}

		response.Tasks = append(response.Tasks, taskResponse)
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// GetAllBoards godoc
// @Summary Get all boards
// @Description Retrieve all boards.  This endpoint requires cookie authentication.
// @Tags boards
// @Accept json
// @Produce json
// @Security CookieAuth
// @Success 200 {array} web.BoardResponse
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /boards [get]
func (c *BoardController) GetAllBoards(ctx *fiber.Ctx) error {
	boards, err := c.boardService.GetAllBoards()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	response := make([]map[string]interface{}, len(boards))
	for i, board := range boards {
		boardMap := map[string]interface{}{
			"board_id":   board.ID,
			"name_board": board.NameBoard,
			"tasks":      board.Tasks,
			"board_created_by": map[string]interface{}{
				"user_id":    board.UserID,
				"user_email": board.User.Email,
			},
		}
		response[i] = boardMap
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// DeleteBoardById godoc
// @Summary Delete a board
// @Description Delete a board by its ID.  This endpoint requires cookie authentication.
// @Tags boards
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param   request path int true "Board ID parameter" minimum(1) example(1)
// @Success 200 {object} web.SuccessResponse
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /board/{id} [delete]
func (c *BoardController) DeleteBoardById(ctx *fiber.Ctx) error {
	boardID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid board ID"})
	}

	err = c.boardService.DeleteBoardById(uint64(boardID))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Board deleted successfully"})
}

// EditBoard godoc
// @Summary Edit a board
// @Description Edit a board's name by its ID.  This endpoint requires cookie authentication.
// @Tags boards
// @Accept json
// @Produce json
// @Security CookieAuth
// @Param   request path int true "Board ID parameter" minimum(1) example(1)
// @Param name_board formData string true "New name of the board"
// @Success 200 {object} web.BoardResponse
// @Failure 400 {object} web.ErrorResponse
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /board/{id} [put]
func (c *BoardController) EditBoard(ctx *fiber.Ctx) error {
	boardID, err := ctx.ParamsInt("id")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid board ID"})
	}

	newNameBoard := ctx.FormValue("name_board")
	if newNameBoard == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "New board name is required"})
	}

	updatedBoard, err := c.boardService.EditBoard(uint64(boardID), newNameBoard)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Board updated successfully",
		"board":   updatedBoard,
	})
}
