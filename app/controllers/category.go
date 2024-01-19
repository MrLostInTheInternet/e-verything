package controllers

import (
	"encoding/json"
	"net/http"
	"webapp/app/config"
	"webapp/app/models"

	"github.com/labstack/echo/v4"
)

func SortByCategoryID(c echo.Context) (err error) {
	var products []models.Product
	// Parse the category ID from the path parameter
	categoryID := c.Param("id")
	result := config.DB.Preload("ProductCategory").Where("category_id=?", categoryID).Find(&products)
	if result.Error != nil {
		return c.String(http.StatusBadRequest, "Category not found")
	}
	// Marshal the data with indentation
	formatJSON, err := json.MarshalIndent(products, "", "  ")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}
	// Set content type to JSON and send the pretty-printed JSON as a raw response
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	return c.String(http.StatusOK, string(formatJSON))
}
