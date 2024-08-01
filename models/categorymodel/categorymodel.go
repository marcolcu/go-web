package categorymodel

import (
	"go-web-native/config"
	"go-web-native/entities"
)

func GetAll() []entities.Category {
	var categories []entities.Category
	result := config.DB.Find(&categories)
	if result.Error != nil {
		panic(result.Error)
	}
	return categories
}

func Create(category entities.Category) bool {
	result := config.DB.Create(&category)
	if result.Error != nil {
		panic(result.Error)
	}
	return result.RowsAffected > 0
}

func Detail(id int) entities.Category {
	var category entities.Category
	result := config.DB.First(&category, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return category
}

func Update(id int, category entities.Category) bool {
	var existingCategory entities.Category
	result := config.DB.First(&existingCategory, id)
	if result.Error != nil {
		panic(result.Error)
	}
	existingCategory.Name = category.Name
	existingCategory.UpdatedAt = category.UpdatedAt
	updateResult := config.DB.Save(&existingCategory)
	if updateResult.Error != nil {
		panic(updateResult.Error)
	}
	return updateResult.RowsAffected > 0
}

func Delete(id int) error {
	result := config.DB.Delete(&entities.Category{}, id)
	return result.Error
}
