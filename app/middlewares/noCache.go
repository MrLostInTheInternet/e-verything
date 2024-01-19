package middlewares

import "github.com/labstack/echo/v4"

func NoCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Set headers to instruct the browser not to cache the response
		c.Response().Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0")
		c.Response().Header().Set("Pragma", "no-cache")
		c.Response().Header().Set("Expires", "Fri, 01 Jan 1990 00:00:00 GMT")
		return next(c)
	}
}
