package middleware

import (
	"manajemen_tugas_master/service"

	"github.com/gofiber/fiber/v2"
)

func AuthUser(userService service.UserService) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Get the cookie from request
		tokenStringJwt := ctx.Cookies("Authorization")
		// tokenStringOauth := ctx.Cookies("GoogleAuthorization")

		// validate tokenStringJwt
		if tokenStringJwt == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Your session has expired, Please login again"})
		}
		if tokenStringJwt != "" {
			// Decode and validate the token
			user, err := userService.RequireAuthUser(tokenStringJwt)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
			}

			// menyimpan data user agar bisa di akses jika diperlukan. Dan perlu di ingat data user akan berubah menjadi interface{}, bukan *domain.user lagi.
			ctx.Locals("user", user)
		}

		// // validate tokenStringOauth
		// if tokenStringOauth == "" {
		// 	return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Your session has expired, Please login again"})
		// }
		// if tokenStringOauth != "" {
		// 	email, ok := ctx.Locals("sessionEmail").(string)
		// 	if !ok {
		// 		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Email not found in request context"})
		// 	}
		// 	user, err := userService.RequireOauth(email)
		// 	if err != nil {
		// 		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err})
		// 	}

		// 	// menyimpan data user agar bisa di akses jika diperlukan. Dan perlu di ingat data user akan berubah menjadi interface{}, bukan *domain.user lagi.
		// 	ctx.Locals("userOauth", user)
		// }

		// agar middleware terdapat pada route di bawahnya dan akan terus di eksekusi terlebih dahulu sebelum route di bawahnya.
		return ctx.Next()
	}
}
