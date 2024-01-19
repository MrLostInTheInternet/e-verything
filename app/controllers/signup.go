package controllers

import (
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"
	"webapp/app/config"
	"webapp/app/middlewares"
	"webapp/app/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SignupGET(c echo.Context) (err error) {
	return c.Render(http.StatusOK, "signup", echo.Map{
		"Date": time.Now().Year(),
	})
}

func SignupPOST(c echo.Context) (err error) {
	u := new(models.User)
	if err := c.Bind(u); err != nil {
		return c.String(http.StatusBadRequest, "Failed to signup")
	}
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(u.Email)
	u.Password = strings.TrimSpace(u.Password)
	errorMessages := make(map[string]string)
	if u.Username == "" {
		errorMessages["Username"] = "Username is required"
	}
	if u.Email == "" || !isValidEmail(u.Email) {
		errorMessages["Email"] = "Email is required"
	}
	if u.Password == "" {
		errorMessages["Password"] = "Password is required"
	}
	if len(errorMessages) > 0 {
		return c.Render(http.StatusOK, "signup", echo.Map{
			"Errors": errorMessages,
			"Date":   time.Now().Year(),
		})
	}
	// Server processing
	// Check if the email already exists in the database
	var existingUser models.User
	found := config.DB.Where("email = ?", u.Email).First(&existingUser)
	// If the email is found, return an error
	if found.Error == nil { // No error means a record was found
		return c.JSON(http.StatusConflict, echo.Map{"message": "You already have an account with this email"})
	}
	// If the error is not a 'record not found' error, handle it as a server error
	if !errors.Is(found.Error, gorm.ErrRecordNotFound) {
		// Log the error and return a server error response
		return c.String(http.StatusInternalServerError, "Server error")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to encrypt password")
	}
	u.Password = string(hashedPassword)
	// Save the user credentials in the database
	result := config.DB.Save(u)
	if result.Error != nil {
		return c.String(http.StatusInternalServerError, "Failed to save the user credentials")
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
	return c.Redirect(http.StatusFound, "login")
}

func isValidEmail(email string) bool {
	// Regular expression for validating an email
	emailRegex := regexp.MustCompile(`(?i)^(?:[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_` + "`" + `{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\])$`)
	return emailRegex.MatchString(email)
}

/*
func generateVerificationToken() (string, time.Time, error) {
	tokenBytes := make([]byte, 16) // 128-bit token
	_, err := rand.Read(tokenBytes)
	if err != nil {
		// Log the error
		log.Printf("Error generating verification token: %v", err)
		// Return a token generation error
		return "", time.Time{}, err
	}
	token := hex.EncodeToString(tokenBytes)
	expiry := time.Now().Add(24 * time.Hour) // Token expires after 24 hours
	return token, expiry, nil
}
*/
