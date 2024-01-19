package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func LogoutPOST(c echo.Context) (err error) {
	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.MaxAge = -1
	c.SetCookie(cookie)
	return c.Redirect(http.StatusFound, "login")
}
