package fcategorycontroller

import (
	"github.com/gofiber/fiber/v2"
)

// Index returns a list of categories as JSON
func Index(c *fiber.Ctx) error {
	return c.Render("views/category/index.html", nil)
}

// Add creates a new category from JSON data
func Add(c *fiber.Ctx) error {
	return c.Render("views/category/create.html", nil)
}

// Edit updates a category with JSON data
func Edit(c *fiber.Ctx) error {
	return c.Render("views/category/edit.html", nil)
}
