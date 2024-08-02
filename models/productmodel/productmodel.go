package productmodel

import (
	"go-web-native/config"
	"go-web-native/entities"
	"log"
	"time"
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
        log.Printf("Error fetching product: %v", result.Error)
        return false
    }

    // Ensure that the CategoryId exists in the categories table
    var category entities.Category
    if err := config.DB.First(&category, product.CategoryId).Error; err != nil {
        log.Printf("Category ID %d does not exist: %v", product.CategoryId, err)
        return false
    }

    existingProduct.Name = product.Name
    existingProduct.CategoryId = product.CategoryId
    existingProduct.Stock = product.Stock
    existingProduct.Description = product.Description
    existingProduct.UpdatedAt = time.Now()

    updateResult := config.DB.Save(&existingProduct)
    if updateResult.Error != nil {
        log.Printf("Error updating product: %v", updateResult.Error)
        return false
    }
    return updateResult.RowsAffected > 0
}

func Delete(id int) error {
	result := config.DB.Delete(&entities.Product{}, id)
	return result.Error
}
