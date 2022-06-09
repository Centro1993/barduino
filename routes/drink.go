package routes

import (
	"barduino/logic"
	"barduino/models"

	"github.com/gofiber/fiber/v2"
)

func DrinkRoutes (app *fiber.App) {
	app.Post("/drink", func(c *fiber.Ctx) error {
		drink := new(models.Drink)
	
		if err := c.BodyParser(drink); err != nil {
			return c.Status(503).SendString(err.Error())
		}
		
		if tx := models.DB.Omit("Recipe").Omit("Served").Create(&drink); tx.Error != nil {
			return c.Status(503).SendString(tx.Error.Error())
		}

		// Start Pouring
		if err := logic.PourDrink(*drink); err != nil {
			return c.Status(404).JSON(err.Error())
		}

		return c.Status(201).JSON(drink)
	})

    app.Delete("/drink", func(c *fiber.Ctx) error {
		if err := logic.CancelDrink(); err != nil {
			return c.Status(404).JSON(err.Error())
		}

		return c.SendStatus(201)
	})
}