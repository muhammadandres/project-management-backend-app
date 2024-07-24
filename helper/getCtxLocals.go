package helper

import (
	"errors"
	"manajemen_tugas_master/model/domain"

	"github.com/gofiber/fiber/v2"
)

func GetCtxLocals(ctx *fiber.Ctx) (uint64, error) {
	user := ctx.Locals("user")
	if user != nil {
		if u, ok := user.(*domain.User); ok {
			return u.ID, nil
		}
	}

	userOauth := ctx.Locals("userOauth")
	if userOauth != nil {
		if u, ok := userOauth.(*domain.User); ok {
			return u.ID, nil
		}
	}

	return 0, errors.New("User not found or invalid type")
}
