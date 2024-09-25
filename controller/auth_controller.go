package controller

import (
	"time"

	"go-fiber/model"

	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

type AuthRequest struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	Created_at string `json:"created_at"`
	Updated_at string `json:"updated_at"`
	Deleted_at string `json:"deleted_at"`
}

type AuthController struct {
	DB *gorm.DB
}

func (ac *AuthController) Register(ctx *fiber.Ctx) error {
	var registerRequest AuthRequest

	// Parse the request body into the AuthRequest struct.
	if err := ctx.BodyParser(&registerRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Cannot parse JSON",
		})
	}

	if registerRequest.Name == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Name is required",
		})
	}

	if registerRequest.Email == "" || registerRequest.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Email and password are required",
		})
	}

	// Create a User model instance and set its fields.
	user := model.Users{
		Name:  registerRequest.Name,
		Email: registerRequest.Email,
	}

	// Hash the user's password.
	if err := user.HashPassword(registerRequest.Password); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to hash password",
		})
	}

	// Save the user in the database.
	if err := ac.DB.Create(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to register user",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "User registered successfully",
		"data":    user,
	})
}

var jwtSecret = []byte("your_secret_key")

func (lc *AuthController) Login(ctx *fiber.Ctx) error {
	var loginRequest AuthRequest

	if err := ctx.BodyParser(&loginRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "cannot parse JSON",
		})
	}

	if loginRequest.Email == "" || loginRequest.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Email and password are required",
		})
	}

	var user model.Users
	if err := lc.DB.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		// If the user is not found, return a 401 Unauthorized status.
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Invalid username or password",
		})
	}

	// Check if the provided password matches the hashed password.
	if err := user.CheckPassword(loginRequest.Password); err != nil {
		// If the password does not match, return a 401 Unauthorized status.
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Invalid username or password",
		})
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": loginRequest.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(), // Token expires after 72 hours.
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Could not create token",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "Login Succesfully",
		"token":   tokenString,
	})
}
