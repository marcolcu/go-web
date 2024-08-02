package productcontroller

import (
	"go-web-native/entities"
	"go-web-native/models/categorymodel"
	"go-web-native/models/productmodel"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Index(c *fiber.Ctx) error {
	// Fetch query parameters
	draw := c.Query("draw")
	start := c.Query("start")
	length := c.Query("length")
	search := c.Query("search[value]") // Fetch the search query

	// Parse query parameters for pagination
	startInt, err := strconv.Atoi(start)
	if err != nil {
		startInt = 0
	}
	lengthInt, err := strconv.Atoi(length)
	if err != nil {
		lengthInt = 10
	}

	// Fetch all products
	products := productmodel.GetAll()
	totalRecords := len(products) // Total number of records

	// Filter products based on the search query
	filteredProducts := []entities.Product{}
	for _, product := range products {
		if search == "" || strings.Contains(strings.ToLower(product.Name), strings.ToLower(search)) {
			filteredProducts = append(filteredProducts, product)
		}
	}
	filteredRecords := len(filteredProducts)

	// Slice products according to pagination parameters
	end := startInt + lengthInt
	if end > filteredRecords {
		end = filteredRecords
	}
	filteredProducts = filteredProducts[startInt:end]

	// Create DataTables response format
	response := map[string]interface{}{
		"draw":            draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            filteredProducts,
	}

	// Set content type and encode response
	return c.JSON(response)
}
func Detail(c *fiber.Ctx) error {
	// Handle only GET requests
	if c.Method() != fiber.MethodGet {
		return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{"error": "Method not allowed"})
	}

	// Parse the product ID from the URL parameters
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	// Fetch the product details from the database
	product := productmodel.Detail(id)
	if (product == entities.Product{}) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	// Return the product details
	return c.JSON(product)
}

func Add(c *fiber.Ctx) error {
	// Handle only POST requests
	if c.Method() != fiber.MethodPost {
		return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{"error": "Method not allowed"})
	}

	// Parse the JSON payload into the Product struct
	var product entities.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Set timestamps
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	// Create the product in the database
	if ok := productmodel.Create(product); !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
	}

	// Return the created product
	return c.Status(fiber.StatusCreated).JSON(product)
}
func Edit(c *fiber.Ctx) error {
	// Handle GET request to fetch product details by ID
	if c.Method() == fiber.MethodGet {
		id, err := strconv.Atoi(c.Params("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
		}

		product := productmodel.Detail(id)
		if (product == entities.Product{}) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}

		categories := categorymodel.GetAll() // Fetch all categories

		data := map[string]interface{}{
			"product":    product,
			"categories": categories,
		}

		return c.JSON(data)
	}

	// Handle POST request to update product details
	if c.Method() == fiber.MethodPost {
		var product entities.Product

		// Parse product data from the request body
		if err := c.BodyParser(&product); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
		}

		// Extract ID from query parameters
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
		}

		product.Id = uint(id)
		product.UpdatedAt = time.Now()

		// Update the product
		if ok := productmodel.Update(id, product); !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product"})
		}

		return c.Status(fiber.StatusOK).JSON(product)
	}

	// Handle unsupported methods
	return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{"error": "Method not allowed"})
}

func Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product ID"})
	}

	if err := productmodel.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete product"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Product deleted successfully"})
}