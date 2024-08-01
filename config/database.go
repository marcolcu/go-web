package config

import (
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"go-web-native/entities"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "root:indotim@tcp(127.0.0.1:3306)/go_products?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&entities.Product{})
	if err != nil {
		panic("Failed to migrate database schema")
	}
	err = db.AutoMigrate(&entities.Category{})
	if err != nil {
		panic("Failed to migrate database schema")
	}
	err = db.AutoMigrate(&entities.User{})
	if err != nil {
		panic("Failed to migrate database schema")
	}

	log.Println("Database connected and schema migrated")
	DB = db
}
