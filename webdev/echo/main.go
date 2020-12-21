package main

import (
	"crypto/subtle"
	"net/http"
    "github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func authValidator(username, password string, c echo.Context) (bool, error) {
	// Be careful to use constant time comparison to prevent timing attacks
	if subtle.ConstantTimeCompare([]byte(username), []byte("admin")) == 1 &&
		subtle.ConstantTimeCompare([]byte(password), []byte("admin")) == 1 {
		return true, nil
	}
	return false, nil
}

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Skipper: middleware.DefaultSkipper,
		Validator:authValidator,
	}))

	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup: "header:X-CSRF-TOKEN",
	}))

	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "",
		ContentTypeNosniff:    "",
		XFrameOptions:         "",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})



	// Start server
	e.Logger.Fatal(e.Start(":1234"))
}
