package fcategorycontroller

import (
	"net/http"
	"text/template"
)

// Index returns a list of categories as JSON
func Index(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("views/category/index.html")
	if err != nil {
		panic(err)
	}

	temp.Execute(w, nil)
}

// Add creates a new category from JSON data
func Add(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("views/category/create.html")
	if err != nil {
		panic(err)
	}

	temp.Execute(w, nil)
}

// Edit updates a category with JSON data
func Edit(w http.ResponseWriter, r *http.Request) {
	temp, err := template.ParseFiles("views/category/edit.html")
	if err != nil {
		panic(err)
	}

	temp.Execute(w, nil)
}
