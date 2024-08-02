package main

import (
	"go-web-native/config"
	"go-web-native/controllers/authcontroller"
	"go-web-native/controllers/categorycontroller"
	"go-web-native/controllers/frontend/fauthcontroller"
	"go-web-native/controllers/frontend/fcategorycontroller"
	"go-web-native/controllers/frontend/fproductcontroller"
	"go-web-native/controllers/homecontroller"
	"go-web-native/controllers/productcontroller"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectDB()

	app := fiber.New()

	// Redirect root URL to /login
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/login")
	})

	// Auth routes
	app.Post("/login", authcontroller.Login)
	app.Get("/login", authcontroller.Login)
	app.Get("/register", fauthcontroller.Index)
	app.Post("/api/register", authcontroller.Register)
	app.Post("/api/logout", authcontroller.Logout)

	// Middleware to protect routes
	app.Use("/home", authcontroller.AuthMiddleware)
	app.Get("/home", homecontroller.Welcome)

	app.Use("/categories", authcontroller.AuthMiddleware)
	app.Get("/categories", fcategorycontroller.Index)
	app.Get("/categories/add", fcategorycontroller.Add)
	app.Get("/categories/edit", fcategorycontroller.Edit)

	app.Use("/products", authcontroller.AuthMiddleware)
	app.Get("/products", fproductcontroller.Index)
	app.Get("/products/add", fproductcontroller.Add)
	app.Get("/products/edit", fproductcontroller.Edit)
	app.Get("/products/detail", fproductcontroller.Detail)

	app.Use("/api/categories", authcontroller.AuthMiddleware)
	app.Get("/api/categories", categorycontroller.Index)
	app.Post("/api/categories/add", categorycontroller.Add)
	app.Get("/api/categories/edit", categorycontroller.Edit)
	app.Post("/api/categories/edit", categorycontroller.Edit)
	app.Delete("/api/categories/delete", categorycontroller.Delete)

	app.Use("/api/products", authcontroller.AuthMiddleware)
	app.Get("/api/products", productcontroller.Index)
	app.Post("/api/products/add", productcontroller.Add)
	app.Get("/api/products/detail", productcontroller.Detail)
	app.Post("/api/products/edit", productcontroller.Edit)
	app.Delete("/api/products/delete", productcontroller.Delete)

	log.Println("Server running on port 8080")
	log.Fatal(app.Listen(":8080"))
}
