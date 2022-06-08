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
	
        
		
		if tx := models.DB.Omit("Recipe").Create(&drink); tx.Error != nil {
			return c.Status(503).SendString(tx.Error.Error())
		}

		return c.Status(201).JSON(drink)
	})
}