package controller

import (
	"encoding/json"
	"manajemen_tugas_master/model/domain"
	"manajemen_tugas_master/service"
	"strconv"
	"time"

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

	// Ambil cookie yang ada
	existingCookie := ctx.Cookies("BoardIDs")
	var boardIDs []uint64

	if existingCookie != "" {
		// Jika cookie sudah ada, parse nilai yang ada
		err = json.Unmarshal([]byte(existingCookie), &boardIDs)
		if err != nil {
			boardIDs = []uint64{}
		}
	}

	// Tambahkan ID board baru ke array
	boardIDs = append(boardIDs, boardDb.ID)

	// Konversi array kembali ke JSON
	newCookieValue, err := json.Marshal(boardIDs)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update cookie"})
	}

	// Set cookie dengan daftar board IDs
	ctx.Cookie(&fiber.Cookie{
		Name:    "BoardIDs",
		Value:   string(newCookieValue),
		Expires: time.Now().AddDate(100, 0, 0), // 100 tahun
	})

	// Set cookie untuk board ID terbaru
	ctx.Cookie(&fiber.Cookie{
		Name:    "LatestBoardID",
		Value:   strconv.FormatUint(uint64(boardDb.ID), 10),
		Expires: time.Now().AddDate(100, 0, 0), // 100 tahun
	})

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Board created successfully", "board_id": boardDb.ID})
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

	// Remove the deleted board ID from the cookie
	existingCookie := ctx.Cookies("BoardIDs")
	var boardIDs []uint64

	if existingCookie != "" {
		err = json.Unmarshal([]byte(existingCookie), &boardIDs)
		if err == nil {
			// Remove the deleted board ID
			for i, id := range boardIDs {
				if id == uint64(boardID) {
					boardIDs = append(boardIDs[:i], boardIDs[i+1:]...)
					break
				}
			}

			// Update the cookie
			newCookieValue, err := json.Marshal(boardIDs)
			if err == nil {
				ctx.Cookie(&fiber.Cookie{
					Name:    "BoardIDs",
					Value:   string(newCookieValue),
					Expires: time.Now().AddDate(100, 0, 0), // 100 years
				})
			}
		}
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
