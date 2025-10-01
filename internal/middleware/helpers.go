package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthHelpers struct{}

func NewAuthHelpers() *AuthHelpers {
	return &AuthHelpers{}
}

func (a *AuthHelpers) SetAccessTokenCookie(c *fiber.Ctx, accessToken string) {
	secure := "development" == "production"
	sameSite := "Lax"

	if secure {
		sameSite = "None"
	}

	c.Cookie(&fiber.Cookie{
		Name:     "vazzaccess",
		Value:    accessToken,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		Domain:   "",
		Secure:   secure,
		HTTPOnly: true,
		SameSite: sameSite,
	})
}

func (a *AuthHelpers) ClearAuthCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     "vazzaccess",
		Value:    "",
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HTTPOnly: true,
	})
}

func (a *AuthHelpers) GetTokenFromCookie(c *fiber.Ctx) string {
	return c.Cookies("vazzaccess")
}

func (a *AuthHelpers) GetTokenFromHeader(c *fiber.Ctx) string {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return ""
	}

	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}

	return ""
}

func (a *AuthHelpers) GetToken(c *fiber.Ctx) string {
	if token := a.GetTokenFromCookie(c); token != "" {
		return token
	}

	return a.GetTokenFromHeader(c)
}
