package controller

import (
	"context"
	"encoding/json"
	"go-fiber/database"
	"go-fiber/model"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type RedisController struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func (uc *RedisController) GetUsers(ctx *fiber.Ctx) error {
	val, err := database.GetKey(ctx.Context(), "users")
	if err == redis.Nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"message": "Users not found",
		})
	} else if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to connect to Redis",
		})
	}

	var users []model.UsersRedis
	if err := json.Unmarshal([]byte(val), &users); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to decode user data",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success",
		"data":    users,
	})
}

func (uc *RedisController) SyncUserToRedisSync() error {
	var users []model.UsersRedis

	if err := uc.DB.Model(&model.Users{}).Select("ID", "Name", "Email").Scan(&users).Error; err != nil {
		return err
	}

	userJSON, err := json.Marshal(users)
	if err != nil {
		return err
	}

	if err := database.SetKey(context.Background(), "users", userJSON); err != nil {
		return err
	}

	return nil
}

// SyncUserToRedis handles the HTTP request
func (uc *RedisController) SyncUserToRedis(c *fiber.Ctx) error {
	if err := uc.SyncUserToRedisSync(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to synchronize user data.",
		})
	}

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "User data synchronized successfully.",
	})
}

func (uc *RedisController) SetUserToRedis(c *fiber.Ctx) error {
	userID := c.Params("id")

	var user model.UsersRedis

	if err := uc.DB.Model(&model.Users{}).Select("ID", "Name", "Email").First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  fiber.StatusNotFound,
				"message": "User not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to fetch user",
		})
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to encode user data",
		})
	}

	err = database.SetKey(c.Context(), "profile", userJSON)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  fiber.StatusInternalServerError,
			"message": "Failed to connect to Redis",
		})
	}

	return c.JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"message": "success",
	})
}
