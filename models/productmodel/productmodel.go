package productmodel

import (
	"go-web-native/config"
	"go-web-native/entities"
)

func GetAll() []entities.Product {
	var products []entities.Product
	result := config.DB.Preload("Category").Find(&products)
	if result.Error != nil {
		panic(result.Error)
	}
	return products
}

func Create(product entities.Product) bool {
	result := config.DB.Create(&product)
	if result.Error != nil {
		panic(result.Error)
	}
	return result.RowsAffected > 0
}

func Detail(id int) entities.Product {
	var product entities.Product
	result := config.DB.Preload("Category").First(&product, id)
	if result.Error != nil {
		panic(result.Error)
	}
	return product
}

func Update(id int, product entities.Product) bool {
	var existingProduct entities.Product
	result := config.DB.First(&existingProduct, id)
	if result.Error != nil {
		panic(result.Error)
	}
	existingProduct.Name = product.Name
	existingProduct.Category.Id = product.Category.Id
	existingProduct.Stock = product.Stock
	existingProduct.Description = product.Description
	existingProduct.UpdatedAt = product.UpdatedAt
	updateResult := config.DB.Save(&existingProduct)
	if updateResult.Error != nil {
		panic(updateResult.Error)
	}
	return updateResult.RowsAffected > 0
}

func Delete(id int) error {
	result := config.DB.Delete(&entities.Product{}, id)
	return result.Error
}
