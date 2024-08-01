package categorycontroller

import (
	"encoding/json"
	"go-web-native/entities"
	"go-web-native/models/categorymodel"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Index returns a list of categories as JSON
func Index(w http.ResponseWriter, r *http.Request) {
	// Fetch query parameters
	draw := r.URL.Query().Get("draw")
	start := r.URL.Query().Get("start")
	length := r.URL.Query().Get("length")
	search := r.URL.Query().Get("search[value]") // Fetch the search query

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
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

// Add creates a new category from JSON data
func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var category entities.Category

		err := json.NewDecoder(r.Body).Decode(&category)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		category.CreatedAt = time.Now()
		category.UpdatedAt = time.Now()

		if ok := categorymodel.Create(category); !ok {
			http.Error(w, "Failed to create category", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(category)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// Edit updates a category with JSON data
func Edit(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var category entities.Category

		err := json.NewDecoder(r.Body).Decode(&category)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		idString := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		category.UpdatedAt = time.Now()

		if ok := categorymodel.Update(id, category); !ok {
			http.Error(w, "Failed to update category", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(category)
		return
	} else if r.Method == http.MethodGet {
		idString := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		category, err := categorymodel.Detail(id)
		if err != nil {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(category)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// Delete removes a category by ID
func Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		idString := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Invalid category ID", http.StatusBadRequest)
			return
		}

		if err := categorymodel.Delete(id); err != nil {
			http.Error(w, "Failed to delete category", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
