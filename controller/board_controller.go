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

func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	board := &domain.Board{}
	nameBoard := ctx.FormValue("name_board")
	board.NameBoard = nameBoard

	boardDb, err := c.boardService.CreateBoard(board)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Board berhasil dibuat",
		"board_id": boardDb.ID,
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

	return ctx.Status(fiber.StatusOK).JSON(board)
}

func (c *BoardController) GetAllBoards(ctx *fiber.Ctx) error {
	boards, err := c.boardService.GetAllBoards()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(boards)
}

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
