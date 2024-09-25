package controller

import (
	"go-fiber/model"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UserController struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func (uc *UserController) Profile(ctx *fiber.Ctx) error {
	// Get the user ID from the context and ensure it's an int
	userID, ok := ctx.Locals("userID").(int)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "Invalid user ID",
		})
	}

	var user model.Users
	if err := uc.DB.First(&user, userID).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "User not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success",
		"data":    user,
	})
}

func (uc *UserController) GetAllUsers(ctx *fiber.Ctx) error {
	// Default values for pagination
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	// Calculate offset
	offset := (page - 1) * limit

	var users []model.Users
	var total int64

	if err := uc.DB.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Could not fetch users",
		})
	}

	// Count total users
	uc.DB.Model(&model.Users{}).Count(&total)

	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success",
		"data":    users,
		"total":   total,
		"page":    page,
		"limit":   limit,
	})
}

func (uc *UserController) GetUsersByID(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")

	var user model.Users
	if err := uc.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  fiber.StatusNotFound,
				"message": "User not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to fetch user",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success",
		"data":    user,
	})
}

func (uc *UserController) SearchUserByName(ctx *fiber.Ctx) error {
	name := ctx.Query("name")

	var users []model.Users
	if err := uc.DB.Where("name LIKE ?", "%"+name+"%").Find(&users).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to search users",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success",
		"data":    users,
	})
}

func (uc *UserController) UpdateUser(ctx *fiber.Ctx) error {
	userID := ctx.Params("id")

	var user model.Users
	if err := uc.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  fiber.StatusNotFound,
				"message": "User not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to fetch user",
		})
	}

	var updateRequest model.Users
	if err := ctx.BodyParser(&updateRequest); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "Cannot parse JSON",
		})
	}

	// Check if password needs to be updated
	if updateRequest.Password != "" {
		if err := user.HashPassword(updateRequest.Password); err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  fiber.StatusInternalServerError,
				"message": "Failed to hash password",
			})
		}
		updateRequest.Password = user.Password
	} else {
		// Keep the current password if not provided
		updateRequest.Password = user.Password
	}

	if err := uc.DB.Model(&user).Updates(updateRequest).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to update user",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "User updated successfully",
		"data":    user,
	})
}

func (uc *UserController) DeleteUser(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	var user model.Users
	if err := uc.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  fiber.StatusNotFound,
				"message": "User not found",
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to fetch user",
		})
	}

	if err := uc.DB.Delete(&user).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to delete user",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "User deleted successfully",
	})
}
