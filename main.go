package main

import (
	"go-web-native/config"
	"go-web-native/controllers/authcontroller"
	"go-web-native/controllers/categorycontroller"
	"go-web-native/controllers/frontend/fauthcontroller"
	"go-web-native/controllers/frontend/fcategorycontroller"
	"go-web-native/controllers/homecontroller"
	"go-web-native/controllers/productcontroller"
	"log"
	"net/http"
)

func main() {
	config.ConnectDB()

	// 1. Auth
	http.HandleFunc("/login", authcontroller.Login)

	http.HandleFunc("/register", fauthcontroller.Index)
	http.HandleFunc("/api/register", authcontroller.Register)
	http.HandleFunc("/api/logout", authcontroller.Logout)

	// Middleware to protect routes
	http.Handle("/home", authcontroller.AuthMiddleware(http.HandlerFunc(homecontroller.Welcome)))

	
	http.Handle("/categories", authcontroller.AuthMiddleware(http.HandlerFunc(fcategorycontroller.Index)))
	http.Handle("/categories/add", authcontroller.AuthMiddleware(http.HandlerFunc(fcategorycontroller.Add)))
	http.Handle("/categories/edit", authcontroller.AuthMiddleware(http.HandlerFunc(fcategorycontroller.Edit)))

	http.Handle("/api/categories", authcontroller.AuthMiddleware(http.HandlerFunc(categorycontroller.Index)))
	http.Handle("/api/categories/add", authcontroller.AuthMiddleware(http.HandlerFunc(categorycontroller.Add)))
	http.Handle("/api/categories/edit", authcontroller.AuthMiddleware(http.HandlerFunc(categorycontroller.Edit)))
	http.Handle("/api/categories/delete", authcontroller.AuthMiddleware(http.HandlerFunc(categorycontroller.Delete)))

	http.Handle("/products", authcontroller.AuthMiddleware(http.HandlerFunc(productcontroller.Index)))
	http.Handle("/products/add", authcontroller.AuthMiddleware(http.HandlerFunc(productcontroller.Add)))
	http.Handle("/products/detail", authcontroller.AuthMiddleware(http.HandlerFunc(productcontroller.Detail)))
	http.Handle("/products/edit", authcontroller.AuthMiddleware(http.HandlerFunc(productcontroller.Edit)))
	http.Handle("/products/delete", authcontroller.AuthMiddleware(http.HandlerFunc(productcontroller.Delete)))

	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}
