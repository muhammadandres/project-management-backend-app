package middleware

import (
	"log"
	"manajemen_tugas_master/service"

	"github.com/gofiber/fiber/v2"
)

func RequireAuthUser(userService service.UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Get the cookie from request
		tokenString := ctx.Cookies("Authorization")
		log.Println("tokenstring= " + tokenString)
		if tokenString == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Your session has expired, Please login again"})
		}

		// Decode and validate the token
		user, err := userService.RequireAuthUser(tokenString)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Please log in to access this menu"})
		}

		// menyimpan data user agar bisa di akses jika diperlukan. Dan perlu di ingat data user akan berubah menjadi interface{}, bukan *domain.user lagi.
		ctx.Locals("user", user)

		// agar middleware terdapat pada route di bawahnya dan akan terus di eksekusi terlebih dahulu sebelum route di bawahnya.
		return ctx.Next()
	}
}
