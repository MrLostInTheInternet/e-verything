package router

import (
	"webapp/app/controllers/admin"
	"webapp/app/controllers/authentication"
	"webapp/app/controllers/shop"
	"webapp/app/middlewares"

	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) {
	// Management endpoint
	e.GET("/management", admin.ManagementGET, middlewares.WithAuth)

	// Signup and Login endpoints
	e.GET("/signup", authentication.SignupGET)
	e.POST("/signup", authentication.SignupPOST)
	e.GET("/login", authentication.LoginGET)
	e.POST("/login", authentication.LoginPOST)

	// Logout endpoint
	e.POST("/logout", authentication.LogoutPOST, middlewares.WithAuth)

	// Products endpoints

	// Get all products
	e.GET("/api/products", shop.GetAllProducts, middlewares.WithAuth)
	// Get product by id
	e.GET("/api/products/:id", shop.GetProductByID, middlewares.WithAuth)
	// Update product by id
	e.PUT("/api/products/:id", shop.UpdateProductByID, middlewares.WithAuth)
	// Get product by id to edit in table row
	e.GET("/edit/:id", admin.GetProductToEdit, middlewares.WithAuth)
	// Add new product
	e.POST("/api/products", shop.CreateProduct, middlewares.WithAuth)
	// Delete product by id
	e.DELETE("/api/products/:id", shop.DeleteProductByID, middlewares.WithAuth)

	// Category endpoints
	e.GET("/api/products/c/:id", shop.SortByCategoryID, middlewares.WithAuth)

	// Dashboard
	e.GET("/dashboard", shop.DashboardGET, middlewares.WithAuth)
	e.GET("/search", shop.SearchHandler, middlewares.WithAuth)
	// e.GET("/sort", shop.SortHandler, middlewares.WithAuth)
	// Product details
	e.GET("/details/:id", shop.GetProductDetails, middlewares.WithAuth)
	e.GET("/products", shop.GetListProducts, middlewares.WithAuth)
}
