package main

import (
	"barduino/models"
	"barduino/routes"
	"barduino/gpio"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2"
	"fmt"
)

func main() {
	Setup()
}

func Setup() *fiber.App {
	app := fiber.New()


	app.Use(logger.New(logger.Config{
		Format: fmt.Sprintf("[${time}] method=${method} uri=${path} status=${status} time=${latency}\n"),
	}))
	
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders:  "*",
	}))
	
	app.Static("/", "./Frontend")

	routes.RecipeRoutes(app)
	routes.PumpRoutes(app)
	routes.DrinkRoutes(app)
	routes.StateRoutes(app)

	models.InitDabase()
	if err := gpio.InitGPIO(); err != nil {
		fmt.Println(err.Error())
	}

	app.Listen(":3000")

	return app
}