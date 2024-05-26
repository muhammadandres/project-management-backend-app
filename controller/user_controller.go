package controller

import (
	"github.com/gofiber/fiber/v2"
	"manajemen_tugas_master/model/domain"
	"manajemen_tugas_master/model/web"
	"manajemen_tugas_master/service"
	"strconv"
	"time"
)

type UserController struct {
	userService service.UserService
}

// NewUserController NewUserService menggabungkan UserController dan Userservice untuk membuat instance UserController baru,
// yang memiliki kemampuan UserService
func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService}
}

func (c *UserController) SignupUser(ctx *fiber.Ctx) error {
	var user domain.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	signupUser, err := c.userService.SignupUser(&user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusCreated).JSON(web.CreateResponseUser(signupUser))
}

func (c *UserController) LoginUser(ctx *fiber.Ctx) error {
	var user *domain.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	tokenString, err := c.userService.LoginUser(user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24 * 1),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login successfully"})
}

func (c *UserController) GetUserByID(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	// Konversi userId ke uint64
	userIdUint64, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user Id"})
	}

	user, err := c.userService.GetUserByID(uint(userIdUint64))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return ctx.Status(fiber.StatusOK).JSON(web.CreateResponseUser(user))
}

func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	users, err := c.userService.FindAllUsers()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "No users found"})
	}

	// Parsing semua data user ke struktur WebResponse
	response := make([]web.WebResponse, len(users))
	for i, user := range users {
		response[i] = web.CreateResponseUser(user) // Perhatikan penggunaan *user di sini
	}
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *UserController) UpdateUser(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	// Konversi userId ke uint64
	userIdUint64, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user Id"})
	}

	_, err = c.userService.GetUserByID(uint(userIdUint64))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User not found"})
	}

	var user domain.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	user.ID = userIdUint64

	updateUser, err := c.userService.UpdateUser(&user)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(web.CreateResponseUser(updateUser))
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error {
	userId := ctx.Params("id")
	// Konversi userId ke uint64
	userIdUint64, err := strconv.ParseUint(userId, 10, 64)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user Id"})
	}

	// Cek apakah user dengan Id tersebut ada
	_, err = c.userService.GetUserByID(uint(userIdUint64))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	// Hapus user
	if err := c.userService.DeleteUser(uint(userIdUint64)); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Deleted successfully"})
}
