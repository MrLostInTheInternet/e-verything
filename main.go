package main

import (
	"os"
	"webapp/app/config"
	"webapp/app/middlewares"
	"webapp/app/router"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	config.LoadEnvVariables()
	config.ConnectToDatabase()
}

func main() {
	e := echo.New()

	//e.Use(middleware.CORS())
	//e.Use(middleware.CSRF())
	e.Use(middleware.Logger())
	e.Use(middlewares.NoCache)
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20))) // 20 requests per minute per IP
	e.Static("/public/scripts/htmx.min.js", "/web/public/scripts/htmx.min.js")
	e.Static("/assets", "/web/public/assets")
	config.NewTemplateRenderer(e, "web/public/views/*.html")

	router.Routes(e)

	e.Start(":" + os.Getenv("PORT"))
}
