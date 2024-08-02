package categorycontroller

import (
	"encoding/json"
	"go-web-native/entities"
	"go-web-native/models/categorymodel"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Index returns a list of categories as JSON
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

	categories := categorymodel.GetAll() // Fetch all categories
	totalRecords := len(categories)      // Total number of records

	// Filter categories based on the search query
	filteredCategories := []entities.Category{}
	for _, category := range categories {
		if search == "" || strings.Contains(strings.ToLower(category.Name), strings.ToLower(search)) {
			filteredCategories = append(filteredCategories, category)
		}
	}
	filteredRecords := len(filteredCategories)

	// Slice categories according to pagination parameters
	end := startInt + lengthInt
	if end > filteredRecords {
		end = filteredRecords
	}
	filteredCategories = filteredCategories[startInt:end]

	// Create DataTables response format
	response := map[string]interface{}{
		"draw":            draw,
		"recordsTotal":    totalRecords,
		"recordsFiltered": filteredRecords,
		"data":            filteredCategories,
	}

	// Set content type and encode response
	c.Response().Header.Set("Content-Type", "application/json")
	if err := json.NewEncoder(c.Response().BodyWriter()).Encode(response); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to encode JSON")
	}
	return nil
}

// Add creates a new category from JSON data
func Add(c *fiber.Ctx) error {
	var category entities.Category

	// Parse JSON body into category struct
	if err := c.BodyParser(&category); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
	}

	// Ensure that Name field is not empty
	if category.Name == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Category name is required")
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	// Create the category in the database
	if ok := categorymodel.Create(category); !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to create category")
	}

	// Respond with the created category
	return c.Status(fiber.StatusCreated).JSON(category)
}


// Edit updates a category with JSON data
func Edit(c *fiber.Ctx) error {
	switch c.Method() {
	case fiber.MethodGet:
		// GET request: Retrieve category by ID
		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid category ID")
		}

		category, err := categorymodel.Detail(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).SendString("Category not found")
		}

		return c.JSON(category)

	case fiber.MethodPost:
		// POST request: Update category
		var category entities.Category

		if err := c.BodyParser(&category); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid request payload")
		}

		id, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid category ID")
		}

		if category.Name == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Category name is required")
		}

		category.Id = uint(id)
		category.UpdatedAt = time.Now()

		if ok := categorymodel.Update(id, category); !ok {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to update category")
		}

		return c.Status(fiber.StatusOK).JSON(category)

	default:
		return c.Status(fiber.StatusMethodNotAllowed).SendString("Method not allowed")
	}
}

// Delete removes a category by ID
func Delete(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid category ID")
	}

	if err := categorymodel.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to delete category")
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}
