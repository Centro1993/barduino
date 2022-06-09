package routes

import (
	"github.com/gofiber/fiber/v2"
	"barduino/models"
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

        // calculate recipe serving sizes

        // dont create drinks on empty pumps

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