package admin

import (
	"net/http"
	"webapp/app/config"
	"webapp/app/models"

	"github.com/labstack/echo/v4"
)

func GetProductToEdit(c echo.Context) (err error) {
	productID := c.Param("id")
	// Get the product by id from the database
	var product models.Product
	if err := config.DB.Preload("ProductCategory").First(&product, productID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Product not found"})
	}

	if c.Request().Header.Get("HX-Request") == "true" {
		return c.Render(http.StatusOK, "productEdit", map[string]interface{}{
			"Product": product,
		})
	}

	// Default to returning JSON
	return c.JSON(http.StatusOK, product)
}
