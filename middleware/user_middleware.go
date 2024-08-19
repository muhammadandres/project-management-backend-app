package middleware

import (
	"fmt"
	"manajemen_tugas_master/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func AuthUser(userService service.UserService, store *session.Store) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Get the cookies from request
		// tokenStringJwt := ctx.Cookies("Authorization")
		tokenStringJwt := ctx.Get("Authorization")
		tokenStringOauth := ctx.Get("GoogleAuthorization")
		// tokenStringOauth := ctx.Cookies("GoogleAuthorization")

		fmt.Println("Authorization Token:", tokenStringJwt)
		fmt.Println("GoogleAuthorization Token:", tokenStringOauth)

		// Validate tokenStringJwt
		if tokenStringJwt != "" {
			// Decode and validate the token
			user, err := userService.RequireAuthUser(tokenStringJwt)
			if err != nil {
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid JWT token"})
			}

			// Store user data in context
			ctx.Locals("user", user)
		}

		// Validate tokenStringOauth
		if tokenStringOauth != "" {
			sess, err := store.Get(ctx)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mendapatkan sesi"})
			}

			email := sess.Get("email")
			if email == nil {
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Google OAuth session not found"})
			}

			emailStr, ok := email.(string)
			if !ok {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid session email format"})
			}

			user, err := userService.RequireOauth(emailStr)
			if err != nil {
				return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Google OAuth session"})
			}

			ctx.Locals("userOauth", user)
		}

		// If neither JWT nor OAuth token is present, return unauthorized
		if tokenStringJwt == "" && tokenStringOauth == "" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Your session has expired, Please login again"})
		}

		// Continue to the next middleware or route handler
		return ctx.Next()
	}
}
