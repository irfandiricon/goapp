package routes

import (
	"ircn/controller"
	"ircn/middleware"

	"github.com/gofiber/fiber/v2"
)

// RouteConfig holds the route configuration.
type RouteConfig struct {
	App             *fiber.App
	AuthController  *controller.AuthController
	UserController  *controller.UserController
	RedisController *controller.RedisController
}

// Setup sets up the routes.
func (c *RouteConfig) Setup() {
	api := c.App.Group("/api")

	api.Post("/register", c.AuthController.Register)
	api.Post("/login", c.AuthController.Login)

	// Apply AuthMiddleware to protected routes
	protected := api.Group("/", middleware.AuthMiddleware())
	// Add your protected routes here
	protected.Get("/profile", c.UserController.Profile)
	protected.Get("/users/data", c.UserController.GetAllUsers)
	protected.Get("/users/search/:id", c.UserController.GetUsersByID)
	protected.Get("/users/search", c.UserController.SearchUserByName)
	protected.Put("/users/update/:id", c.UserController.UpdateUser)
	protected.Post("/users/create", c.AuthController.Register)
	protected.Delete("/users/delete/:id", c.UserController.DeleteUser)

	// Routes for Redis
	protected.Get("/redis/get-users", c.RedisController.GetUsers)
	protected.Get("/redis/sync-users", c.RedisController.SyncUserToRedis)
	protected.Post("/redis/users/:id", c.RedisController.SetUserToRedis)
}
