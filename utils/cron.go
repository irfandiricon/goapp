package utils

import (
	"fmt"

	"go-fiber/controller"

	"github.com/robfig/cron/v3"
)

var redisController *controller.RedisController

// InitializeCronJob initializes the cron jobs
func InitializeCronJob(redisCtrl *controller.RedisController) {
	redisController = redisCtrl // Set the RedisController to be used in cron jobs

	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc("0 * * * * *", func() {
		/*
			if err := redisController.SyncUserToRedisSync(); err != nil {
				fmt.Println("Error syncing users to Redis every minute:", err)
			}
		*/
	})

	if err != nil {
		fmt.Println("Error adding cron job every minute:", err)
	}

	_, err = c.AddFunc("0 */5 * * * *", func() {
		/*
			if err := redisController.SyncUserToRedisSync(); err != nil {
				fmt.Println("Error syncing users to Redis every 2 minutes:", err)
			}
		*/
	})
	if err != nil {
		fmt.Println("Error adding cron job every 2 minutes:", err)
	}

	c.Start()

	fmt.Println("Connected to Cron!")
}
