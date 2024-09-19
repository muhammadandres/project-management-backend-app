package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"manajemen_tugas_master/helper"
	"manajemen_tugas_master/model/domain"
	"manajemen_tugas_master/model/web"
	"manajemen_tugas_master/service"
	"net/url"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/oauth2"
)

type UserController struct {
	userService service.UserService
	store       *session.Store
}

func NewUserController(userService service.UserService, store *session.Store) *UserController {
	return &UserController{
		userService: userService,
		store:       store,
	}
}

// SignupUser godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body web.SignupRequest true "User signup information"
// @Success      201  {object}  web.TokenResponse
// @Failure      400  {object}  web.ErrorResponse
// @Failure      409  {object}  web.ErrorResponse
// @Failure      500  {object}  web.ErrorResponse
// @Header       200 {string} Set-Cookie "Authorization"
// @Router       /user/signup [post]
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

	// Menentukan domain
	domain := "127.0.0.1"
	if ctx.Hostname() == "manajementugas.com" {
		domain = "manajementugas.com"
	}

	ctx.Cookie(&fiber.Cookie{
		Name:    "Authorization",
		Path:    "/",
		Domain:  domain,
		Value:   tokenString,
		Expires: time.Now().Add(time.Hour * 24 * 3),
	})

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Signup successfully", "token": tokenString})
}

// LoginUser godoc
// @Summary      Authenticate a user
// @Description  Login with user credentials
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body web.LoginRequest true "User login credentials"
// @Success      200  {object}  web.TokenResponse
// @Failure      400  {object}  web.ErrorResponse
// @Failure      401  {object}  web.ErrorResponse
// @Failure      404  {object}  web.ErrorResponse
// @Failure      500  {object}  web.ErrorResponse
// @Header       200 {string} Set-Cookie "Authorization"
// @Router       /user/login [post]
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

	// Menentukan domain
	domain := "127.0.0.1"
	if ctx.Hostname() == "manajementugas.com" {
		domain = "manajementugas.com"
	}

	ctx.Cookie(&fiber.Cookie{
		Name:    "Authorization",
		Path:    "/",
		Domain:  domain,
		Value:   tokenString,
		Expires: time.Now().Add(time.Hour * 24 * 3),
	})

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Signup successfully", "token": tokenString})
}

// GoogleOauth godoc
// @Summary      Initiate Google OAuth
// @Description  Start the Google OAuth process. If successful, the user will be redirected to the URL "(frontendURL)/auth-success?email=(encodedUserEmail)&token=(encodedToken)" with the user's email and token in the query parameters.
// @Tags         users
// @Accept  json
// @Produce      json
// @Success      302  "Redirect to success URL"
// @Failure      400  {object}  web.ErrorResponse
// @Failure      500  {object}  web.ErrorResponse
// @Router       /auth/oauth [get]
func (c *UserController) GoogleOauth(ctx *fiber.Ctx) error {
	config := helper.SetupGoogleAuth()
	url := config.AuthCodeURL("state", oauth2.AccessTypeOffline)

	fmt.Println("Authorization URL:", url)
	fmt.Println("Request headers:", ctx.GetReqHeaders())

	return nil
}

func (c *UserController) GoogleCallback(ctx *fiber.Ctx) error {
	code := ctx.Query("code")

	config := helper.SetupGoogleAuth()
	t, err := config.Exchange(context.Background(), code)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

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

	email, ok := userData["email"].(string)
	if !ok {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Email not found in user data"})
	}

	if err := c.userService.GoogleOauth(email); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Menggunakan sesi untuk menyimpan email
	sess, err := c.store.Get(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mendapatkan sesi"})
	}

	sess.Set("email", email)
	if err := sess.Save(); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan sesi"})
	}

	domain := "127.0.0.1"
	if ctx.Hostname() == "manajementugas.com" {
		domain = "manajementugas.com"
	}

	ctx.Cookie(&fiber.Cookie{
		Name:    "GoogleAuthorization",
		Path:    "/",
		Domain:  domain,
		Value:   t.AccessToken,
		Expires: time.Now().Add(time.Hour * 24 * 3),
	})

	// Redirect ke frontend dengan email sebagai parameter
	frontendURL := "http://127.0.0.1:5173" // Ganti dengan URL frontend

	encodedEmail := url.QueryEscape(email)
	encodedToken := url.QueryEscape(t.AccessToken)
	redirectURL := fmt.Sprintf("%s/auth-success?email=%s&token=%s", frontendURL, encodedEmail, encodedToken)
	return ctx.Redirect(redirectURL)
}

// ForgotPassword godoc
// @Summary Initiate forgot password process
// @Description Send a reset code to the user's email
// @Tags users
// @Accept json
// @Produce json
// @Param        request body web.ForgotPasswordRequest true "User's email"
// @Success      200  {object}  web.SuccessResponse
// @Failure      400  {object}  web.ErrorResponse
// @Failure      500  {object}  web.ErrorResponse
// @Router /user/forgot-password [post]
func (c *UserController) ForgotPassword(ctx *fiber.Ctx) error {
	var request struct {
		Email string `json:"email"`
	}
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err := c.userService.InitiateForgotPassword(request.Email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reset code has been sent to your email",
	})
}

// ResetPassword godoc
// @Summary Reset user's password
// @Description Reset the user's password using the provided reset code
// @Tags users
// @Accept json
// @Produce json
// @Param        request body web.ResetPasswordRequest true "Password reset info"
// @Success      200  {object}  web.SuccessResponse
// @Failure      400  {object}  web.ErrorResponse
// @Failure      500  {object}  web.ErrorResponse
// @Router /user/reset-password [post]
func (c *UserController) ResetPassword(ctx *fiber.Ctx) error {
	var request struct {
		Email       string `json:"email"`
		ResetCode   string `json:"reset_code"`
		NewPassword string `json:"new_password"`
	}
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err := c.userService.ResetPassword(request.Email, request.ResetCode, request.NewPassword)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Password has been reset successfully",
	})
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve user details by user ID
// @Tags users
// @Accept  json
// @Produce  json
// @Param   request path int true "User ID parameter" minimum(1) example(1)
// @Success 200  {object} web.GetUserByIDResponse
// @Failure 400  {object} web.ErrorResponse
// @Failure 404  {object} web.ErrorResponse
// @Failure 500  {object} web.ErrorResponse
// @Router  /users/{id} [get]
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

// GetAllUsers godoc
// @Summary Get all users
// @Description Retrieve information for all users
// @Tags users
// @Accept  json
// @Produce json
// @Response 200 {object} web.GetAllUsersResponse{data=[]web.UserDetail{object,object}}
// @Failure 400  {object} web.ErrorResponse
// @Failure 404  {object} web.ErrorResponse
// @Failure 500  {object} web.ErrorResponse
// @Router /users [get]
func (c *UserController) GetAllUsers(ctx *fiber.Ctx) error {
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

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user's information. This endpoint requires cookie authentication.
// @Tags users
// @Accept  json
// @Produce  json
// @Security CookieAuth
// @Param   request path int true "User ID parameter" minimum(1) example(1)
// @Param        request body web.UpdateUser true "User's email"
// @Success 200  {object} web.UpdateUser
// @Failure 400 {object} web.ErrorResponse
// @Failure 401 {object} web.ErrorResponse "Unauthorized - Cookie authentication required"
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /users/{id} [put]
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

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by their ID. This endpoint requires cookie authentication.
// @Tags users
// @Accept  json
// @Produce  json
// @Security CookieAuth
// @Param   request path int true "User ID parameter" minimum(1) example(1)
// @Success 200  {object} web.SuccessResponse
// @Failure 400 {object} web.ErrorResponse
// @Failure 401 {object} web.ErrorResponse "Unauthorized - Cookie authentication required"
// @Failure 404 {object} web.ErrorResponse
// @Failure 500 {object} web.ErrorResponse
// @Router /users/{id} [delete]
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
