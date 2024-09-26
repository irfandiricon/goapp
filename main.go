package main

import (
	"go-fiber/controller"
	"go-fiber/database"
	"go-fiber/model"
	"go-fiber/routes"
	"go-fiber/utils"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	db := database.ConnectDB()
	redisClient := database.ConnectRedis()

	db.AutoMigrate(&model.Users{})

	authController := &controller.AuthController{DB: db}
	userController := &controller.UserController{DB: db}
	redisController := &controller.RedisController{DB: db, RedisClient: redisClient}
	routeConfig := &routes.RouteConfig{
		App:             app,
		AuthController:  authController,
		UserController:  userController,
		RedisController: redisController,
	}

	utils.InitializeCronJob(redisController)

	routeConfig.Setup()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000" // Default port if not set
	}

	if err := app.Listen(":" + port); err != nil {
		panic(err)
	}
}
