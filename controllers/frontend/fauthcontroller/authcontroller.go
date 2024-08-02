package fauthcontroller

import (
	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	return c.Render("views/auth/register.html", nil)
}
