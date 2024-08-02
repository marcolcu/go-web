package fproductcontroller

import (
	"github.com/gofiber/fiber/v2"
)

// Index returns a list of categories as JSON
func Index(c *fiber.Ctx) error {
	return c.Render("views/product/index.html", nil)
}

// Add creates a new product from JSON data
func Add(c *fiber.Ctx) error {
	return c.Render("views/product/create.html", nil)
}

// Edit updates a product with JSON data
func Edit(c *fiber.Ctx) error {
	return c.Render("views/product/edit.html", nil)
}

// Edit updates a product with JSON data
func Detail(c *fiber.Ctx) error {
	return c.Render("views/product/detail.html", nil)
}
