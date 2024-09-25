package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var JWTSecret = []byte("your_secret_key")

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.StatusUnauthorized,
				"message": "Missing or malformed JWT",
			})
		}

		tokenString := authHeader[len("Bearer "):]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return JWTSecret, nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.StatusUnauthorized,
				"message": "Invalid or expired JWT",
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if userID, ok := claims["user_id"].(float64); ok {
				c.Locals("userID", int(userID))
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"status":  fiber.StatusUnauthorized,
					"message": "Invalid JWT claims",
				})
			}
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  fiber.StatusUnauthorized,
				"message": "Invalid JWT claims",
			})
		}

		return c.Next()
	}
}
