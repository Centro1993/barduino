package routes

import (
	"github.com/gofiber/fiber/v2"
	"barduino/models"
)

func PumpRoutes (app *fiber.App) {
	app.Get("/pump", func(c *fiber.Ctx) error {
		var pumps []models.Pump
	
		if tx := models.DB.Find(&pumps); tx.Error != nil {
            return c.Status(503).SendString(tx.Error.Error())
        }
	
		return c.Status(200).JSON(&pumps)
    })

	app.Post("/pump", func(c *fiber.Ctx) error {
		pump := new(models.Pump)
	
		if err := c.BodyParser(pump); err != nil {
			return c.Status(503).SendString(err.Error())
		}
	
		
		if tx := models.DB.Omit("Pump").Create(&pump); tx.Error != nil {
			return c.Status(503).SendString(tx.Error.Error())
		}

		return c.Status(201).JSON(pump)
	})

	app.Patch("/pump", func(c *fiber.Ctx) error {
        pump := new(models.Pump)
    
        if err := c.BodyParser(pump); err != nil {
            return c.Status(503).SendString(err.Error())
        }
    
        if tx := models.DB.Updates(&pump); tx.Error != nil {
            return c.Status(503).SendString(tx.Error.Error())           
        }
        return c.Status(200).JSON(pump)
	})

    app.Delete("/pump/:id", func(c *fiber.Ctx) error {
        id := c.Params("id")
		
		var ingredient models.Ingredient
		// Delete Ingredients using the Pump first
		if tx := models.DB.Where("pump_id = ?", id).Delete(&ingredient, id); tx.Error != nil {
            return c.Status(503).SendString(tx.Error.Error())           
        }

        var pump models.Pump
        tx := models.DB.Select("Ingredient").Delete(&pump, id)

        if tx.Error != nil {
            return c.Status(503).SendString(tx.Error.Error())           
        }
    
        if tx.RowsAffected == 0 {
            return c.SendStatus(404)
        }
    
        return c.SendStatus(200)
	})
}