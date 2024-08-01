package usermodel

import (
	"go-web-native/config"
	"go-web-native/entities"
)

func CreateUser(user entities.User) error {
	result := config.DB.Create(&user)
	return result.Error
}

func GetUserByEmail(email string) (entities.User, error) {
	var user entities.User
	result := config.DB.Where("email = ?", email).First(&user)
	return user, result.Error
}
