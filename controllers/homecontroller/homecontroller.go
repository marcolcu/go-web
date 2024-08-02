package homecontroller

import (
	"github.com/gofiber/fiber/v2"
)

func Welcome(c *fiber.Ctx) error {
	return c.Render("views/home/index.html", nil)
}