package shop

import (
	"net/http"
	"strings"
	"time"
	"webapp/app/config"
	"webapp/app/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func DashboardGET(c echo.Context) (err error) {
	username := c.Get("username")
	if username == nil {
		return c.String(http.StatusUnauthorized, "User not authenticated")
	}
	usernameStr := username.(string)
	s := cases.Title(language.English)
	return c.Render(http.StatusOK, "dashboard", echo.Map{
		"Username": s.String(usernameStr),
		"Date":     time.Now().Year(),
	})
}

func SearchHandler(c echo.Context) error {
	searchTerm := c.QueryParam("search")

	var products []models.Product
	result := config.DB.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(searchTerm)+"%").Find(&products)
	if result.Error != nil {
		// Handle the error, perhaps return a server error response
		return c.String(http.StatusInternalServerError, result.Error.Error())
	}

	// Directly render your product list as HTML using the Echo context
	return RenderProductListHTML(c, products)
}

// func SortHandler(c echo.Context) error {
// 	sortField := c.QueryParam("sort")

// 	var products []models.Product
// 	config.DB.Order(sortField + " ASC").Find(&products)

// 	htmlContent := RenderProductListHTML(products)
// 	return c.HTML(http.StatusOK, htmlContent)
// }

func RenderProductListHTML(c echo.Context, products []models.Product) error {
	// Generate HTML string for products
	data := map[string]interface{}{
		"Products": products,
	}

	// Render the productList template with the provided data
	return c.Render(http.StatusOK, "productList", data)
}

func GetListProducts(c echo.Context) error {
	var products []models.Product

	// Select * from product
	result := config.DB.Preload("ProductCategory").Find(&products)
	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"error": "Error fetching the products",
		})
	}

	return RenderProductListHTML(c, products)
}
