package router

import (
	"webapp/app/controllers"
	"webapp/app/middlewares"

	"github.com/labstack/echo/v4"
)

func Routes(e *echo.Echo) {
	// management endpoint
	e.GET("/management", controllers.ManagementGET, middlewares.WithAuth)

	// Signup and Login endpoints
	e.GET("/signup", controllers.SignupGET)
	e.POST("/signup", controllers.SignupPOST)
	e.GET("/login", controllers.LoginGET)
	e.POST("/login", controllers.LoginPOST)

	// Logout endpoint
	e.POST("/logout", controllers.LogoutPOST, middlewares.WithAuth)

	// products endpoints
	e.GET("/api/products", controllers.GetAllProducts, middlewares.WithAuth)
	e.GET("/api/products/:id", controllers.GetProductByID, middlewares.WithAuth)
	e.PUT("/api/products/:id", controllers.UpdateProductByID, middlewares.WithAuth)
	e.GET("/api/products/:id/edit", controllers.GetProductToEdit, middlewares.WithAuth)
	e.POST("/api/products", controllers.CreateProduct, middlewares.WithAuth)
	e.DELETE("/api/products/:id", controllers.DeleteProductByID, middlewares.WithAuth)

	// category endpoints
	e.GET("/api/products/c/:id", controllers.SortByCategoryID, middlewares.WithAuth)
}
