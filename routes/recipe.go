package routes

import (
    "github.com/gofiber/fiber/v2"
    "barduino/models"

    "gorm.io/gorm"
)

func RecipeRoutes (app *fiber.App) {
	app.Get("/recipe", func(c *fiber.Ctx) error {
        var recipes []models.Recipe
	
		if tx := models.DB.Preload("Ingredients").Preload("Ingredients.Pump").Find(&recipes); tx.Error != nil {
            return c.Status(503).SendString(tx.Error.Error())
        }
	
		return c.Status(200).JSON(&recipes)
    })

    app.Post("/recipe", func(c *fiber.Ctx) error {
		recipe := new(models.Recipe)
	
		if err := c.BodyParser(recipe); err != nil {
			return c.Status(503).SendString(err.Error())
		}
	
		if tx := models.DB.Omit("Ingredients.Pump").Create(&recipe); tx.Error != nil {
			return c.Status(503).SendString(tx.Error.Error())
		}
		return c.Status(201).JSON(recipe)
	})

    app.Patch("/recipe", func(c *fiber.Ctx) error {
        recipe := new(models.Recipe)
    
        if err := c.BodyParser(recipe); err != nil {
            return c.Status(503).SendString(err.Error())
        }
    
        if tx := models.DB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&recipe); tx.Error != nil {
            return c.Status(503).SendString(tx.Error.Error())           
        }
        return c.Status(200).JSON(recipe)
	})

    app.Delete("/recipe/:id", func(c *fiber.Ctx) error {
        id := c.Params("id")
        var recipe models.Recipe
    
        models.DB.Unscoped().Delete(&recipe, id)

        var ingredients []models.Ingredient
        tx := models.DB.Unscoped().Where("RecipeID = ?", recipe.ID).Delete(&ingredients)

        if tx.Error != nil {
            return c.Status(503).SendString(tx.Error.Error())           
        }
    
        if tx.RowsAffected == 0 {
            return c.SendStatus(404)
        }
    
        return c.SendStatus(200)
	})
}