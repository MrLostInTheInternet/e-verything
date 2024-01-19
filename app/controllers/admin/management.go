package admin

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func ManagementGET(c echo.Context) error {
	username := c.Get("username")
	if username == nil {
		return c.String(http.StatusUnauthorized, "User not authenticated")
	}
	usernameStr := username.(string)
	s := cases.Title(language.English)
	return c.Render(http.StatusOK, "management", echo.Map{
		"Username": s.String(usernameStr),
		"Date":     time.Now().Year(),
	})
}
