package shop

import (
	"webapp/app/config"
	"webapp/app/models"

	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func GetAllProducts(c echo.Context) error {
	// Declare products
	var products []models.Product

	// Check if DB is initialized correctly
	if config.DB == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Database connection is not initialized",
		})
	}

	// Select * from product
	result := config.DB.Preload("ProductCategory").Find(&products)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error fetching the products",
		})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "productsRows", map[string]interface{}{
			"Products": products,
		})
	}

	// Default to returning JSON
	return c.JSON(http.StatusOK, products)
}

func GetProductByID(c echo.Context) error {
	// Declare product
	var product models.Product
	// Check DB init
	if config.DB == nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Database connection is not initialized",
		})
	}
	// Get the product ID from the path parameter
	productID := c.Param("id")
	// Select product_id from product
	result := config.DB.Preload("ProductCategory").First(&product, productID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error fetching the product",
		})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "productRow", map[string]interface{}{
			"Product": product,
		})
	}

	return c.JSON(http.StatusOK, product)
}

func CreateProduct(c echo.Context) (err error) {
	p := new(models.Product)
	// Bind the input form into the model
	if err := c.Bind(p); err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}

	// Update the date that the product has been created
	now := time.Now()
	p.DateCreated = now
	p.LastUpdated = now

	// Create a new product in the database
	result := config.DB.Create(&p)
	if result.Error != nil {
		return c.String(http.StatusBadRequest, "Failed to create new product")
	}

	return c.JSON(http.StatusOK, "Product created successfully")
}

func UpdateProductByID(c echo.Context) (err error) {
	var product models.Product
	input := new(models.ProductInputDTO)
	productID := c.Param("id")
	// Bind the form input into the model input
	if err := c.Bind(input); err != nil {
		return c.String(http.StatusBadRequest, "Bad request")
	}
	// Get the product by id from the database
	if result := config.DB.Preload("ProductCategory").First(&product, productID); result.Error != nil {
		return c.String(http.StatusBadRequest, "Product not found")
	}
	if input.Sku != nil {
		product.Sku = *input.Sku
	}
	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.UnitPrice != nil {
		product.UnitPrice = *input.UnitPrice
	}
	if input.ImageUrl != nil {
		product.ImageUrl = *input.ImageUrl
	}
	if input.UnitsInStock != nil {
		product.UnitsInStock = *input.UnitsInStock
	}
	if input.Active != nil {
		product.Active = *input.Active
	}
	if input.CategoryID != nil {
		product.CategoryID = *input.CategoryID
	}
	product.LastUpdated = time.Now()
	// Save the updated product
	if result := config.DB.Save(&product); result.Error != nil {
		return c.String(http.StatusInternalServerError, "Failed to update product")
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "productRow", map[string]interface{}{
			"Product": product,
		})
	}

	return c.JSON(http.StatusOK, product)
}

func DeleteProductByID(c echo.Context) (err error) {
	var product models.Product
	productID := c.Param("id")
	// Get the product by id from the database
	result := config.DB.Preload("ProductCategory").First(&product, productID)
	if result.Error != nil {
		return c.String(http.StatusBadRequest, "Product not found")
	}
	// Delete the product from the database
	done := config.DB.Delete(&product)
	if done.Error != nil {
		return c.String(http.StatusBadRequest, "Failed to delete product")
	}
	return c.JSON(http.StatusOK, "Product deleted successfully")
}
