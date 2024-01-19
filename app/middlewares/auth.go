package middlewares

import (
	"errors"
	"net/http"
	"os"
	"time"
	"webapp/app/models"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// Authmiddlewares handles JWT authentication
func WithAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the JWT token from the cookie named "tokenName"
		cookie, err := c.Cookie("Authorization")
		if err != nil {
			if err == http.ErrNoCookie {
				// If the cookie is not set, return an unauthorized status
				return c.String(http.StatusUnauthorized, "No token cookie found")
			}
			// For any other type of error, return a bad request status
			return c.String(http.StatusBadRequest, "Error retrieving cookie")
		}
		// Extract the token from the cookie
		tokenString := cookie.Value
		if tokenString == "" {
			return c.String(http.StatusUnauthorized, "Can't get token from cookie")
		}
		token, err := verifyJWTToken(tokenString)
		if err != nil {
			return c.String(http.StatusUnauthorized, "Token is not valid")
		}
		claims := token.Claims.(jwt.MapClaims)
		username := claims["username"]
		// Verify token expiration
		if err := verifyTokenExpiration(token); err != nil {
			return c.String(http.StatusUnauthorized, "Token has expired")
		}
		// You can set the userID in the context for later use in your handlers
		c.Set("username", username)
		return next(c)
	}
}

// GenerateJWTToken generates a new JWT token for the provided user
func GenerateJWTToken(user *models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["userID"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix() // Token expires in 24 hours

	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// VerifyJWTToken verifies the JWT token and returns the token object
func verifyJWTToken(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// VerifyTokenExpiration verifies if the token has expired
func verifyTokenExpiration(token *jwt.Token) error {
	claims := token.Claims.(jwt.MapClaims)
	expirationTime := int64(claims["exp"].(float64))
	currentTime := time.Now().Unix()
	if currentTime > expirationTime {
		return jwt.ValidationError{Inner: errors.New("token has expired")}
	}
	return nil
}
