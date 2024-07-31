package controller

import (
	"context"
	"encoding/json"
	"log"
	"manajemen_tugas_master/helper"
	"manajemen_tugas_master/model/domain"
	"manajemen_tugas_master/model/web"
	"manajemen_tugas_master/service"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
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
	var user *domain.User
	if err := ctx.BodyParser(&user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	tokenString, err := c.userService.SignupUser(user)
	if err != nil {
		switch {
		case err.Error() == "Invalid format in Email":
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
		case err.Error() == "email already exists, please login instead":
			return ctx.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists, please login instead"})
		default:
			// Untuk error lainnya yang tidak terduga
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}
	}

	ctx.Cookie(&fiber.Cookie{
		Name:    "Authorization",
		Value:   tokenString,
		Expires: time.Now().Add(time.Hour * 24 * 3),
	})

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Signup successfully", "token": tokenString})
}

func (c *UserController) LoginUser(ctx *fiber.Ctx) error {
	var user domain.User
	if err := ctx.BodyParser(&user); err != nil {
		log.Printf("Error parsing body: %v", err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}
	log.Printf("Received login request for email: %s, password length: %d", user.Email, len(user.Password))

	tokenString, err := c.userService.LoginUser(&user)
	if err != nil {
		switch err.Error() {
		case "User not found":
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		case "Invalid password":
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid password"})
		case "Invalid format in Email":
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
		default:
			// Untuk error lainnya
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}
	}

	ctx.Cookie(&fiber.Cookie{
		Name:    "Authorization",
		Value:   tokenString,
		Expires: time.Now().Add(time.Hour * 24 * 3),
	})

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Signup successfully", "token": tokenString})
}

func (c *UserController) GoogleOauth(ctx *fiber.Ctx) error {
	config, err := helper.SetupGoogleAuth()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	state := helper.GenerateRandomState()
	url := config.AuthCodeURL(state, oauth2.AccessTypeOffline)
	ctx.Cookie(&fiber.Cookie{
		Name:    "oauthstate",
		Value:   state,
		Expires: time.Now().Add(time.Hour),
	})
	return ctx.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (c *UserController) GoogleCallback(ctx *fiber.Ctx) error {
	state := ctx.Query("state")
	code := ctx.Query("code")

	storedState := ctx.Cookies("oauthstate")
	if state != storedState {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid OAuth state"})
	}

	config, err := helper.SetupGoogleAuth()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	t, err := config.Exchange(context.Background(), code)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	helper.SaveToken(tok)

	client := config.Client(context.Background(), t)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	defer resp.Body.Close()

	var userData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Mendapatkan email dari userData
	email, ok := userData["email"].(string)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Email not found in user data"})
	}

	// validasi user ke repository
	if err := c.userService.GoogleOauth(email); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	ctx.Locals("sessionEmail", email) // Menyimpan email di ctx.locals
	ctx.Cookie(&fiber.Cookie{
		Name:    "GoogleAuthorization",
		Value:   t.AccessToken,
		Expires: time.Now().Add(time.Hour * 24 * 3),
	})

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login successfully"})
}

func (c *UserController) ForgotPassword(ctx *fiber.Ctx) error {
	var request struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err := c.userService.ForgotPassword(request.Email, request.NewPassword)
	if err != nil {
		switch err.Error() {
		case "User not found":
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		case "Invalid format in Email":
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid email format"})
		case "Invalid password format":
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid password format"})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
		}
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password has been reset successfully",
	})
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
	// userEmail := "land45122@gmail.com" // email user
	users, err := c.userService.FindAllUsers()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "No users found"})
	}

	response := make([]web.WebResponse, len(users))
	for i, user := range users {
		response[i] = web.CreateResponseUser(user)
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
