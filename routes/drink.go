package routes

import (
	"barduino/gpio"
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

        // dont create drinks on empty pumps
		if !gpio.CanBeServed(drink.Recipe) {
			c.Status(404)
		}

		// Start Pouring

		return c.Status(201).JSON(drink)
	})

    app.Delete("/drink", func(c *fiber.Ctx) error {
        drink := new(models.Drink)
	
		if err := c.BodyParser(drink); err != nil {
			return c.Status(503).SendString(err.Error())
		}
		
		if tx := models.DB.Omit("Recipe").Omit("Served").Create(&drink); tx.Error != nil {
			return c.Status(503).SendString(tx.Error.Error())
		}

        // calculate recipe serving sizes

		return c.Status(201).JSON(drink)
	})
}