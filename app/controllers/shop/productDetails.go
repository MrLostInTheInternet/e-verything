package shop

import (
	"net/http"
	"webapp/app/config"
	"webapp/app/models"

	"github.com/labstack/echo/v4"
)

func GetProductDetails(c echo.Context) error {
	productID := c.Param("id")
	var product models.Product
	result := config.DB.Preload("ProductCategory").First(&product, productID)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error fetching the product",
		})
	}

	return c.Render(http.StatusOK, "product", echo.Map{
		"Product": product,
	})
}
