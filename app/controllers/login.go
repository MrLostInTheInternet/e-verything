package controllers

import (
	"net/http"
	"time"
	"webapp/app/config"
	"webapp/app/middlewares"
	"webapp/app/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func LoginGET(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "login", echo.Map{
		"Date": time.Now().Year(),
	})
}

func LoginPOST(c echo.Context) (err error) {
	var loginDetails struct {
		Email    string `json:"email" form:"Email"`
		Password string `json:"password" form:"Password"`
	}
	if err := c.Bind(&loginDetails); err != nil {
		return c.String(http.StatusBadRequest, "Invalid input")
	}
	u := new(models.User)
	if result := config.DB.Where("email=?", loginDetails.Email).First(u); result.Error != nil {
		return c.Render(http.StatusUnauthorized, "login", echo.Map{
			"Errors": "Invalid credentials",
			"Date":   time.Now().Year(),
		})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(loginDetails.Password)); err != nil {
		return c.Render(http.StatusUnauthorized, "login", echo.Map{
			"Errors": "Invalid credentials",
			"Date":   time.Now().Year(),
		})
	}
	// Generate JWT token
	token, err := middlewares.GenerateJWTToken(u)
	if err != nil {
		return err
	}
	// Set JWT token as a cookie
	cookie := new(http.Cookie)
	cookie.Name = "Authorization"
	cookie.Value = token
	cookie.HttpOnly = true
	cookie.Expires = time.Now().Add(time.Hour * 24 * 30) // Token expires in 1 month
	c.SetCookie(cookie)
	if c.Request().Header.Get("HX-Request") != "" {
		// It's an HTMX request, instruct the browser to redirect
		return c.HTML(http.StatusOK, `<script>window.location.href = "/management";</script>`)
	} else {
		// Normal request, redirect the server-side way
		return c.Redirect(http.StatusFound, "management")
	}
}
