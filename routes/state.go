package routes


import (
    "github.com/gofiber/fiber/v2"
    "barduino/models"
	"barduino/gpio"
	"barduino/logic"
)

func StateRoutes (app *fiber.App) {
	// Get all Sensor States
	app.Get("/state/sensor", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(&gpio.AverageStateCache)
    })

	app.Get("/state/sensor/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		var pump models.Pump
		if tx := models.DB.First(&pump, id); tx.Error != nil {
            return c.Status(404).SendString(tx.Error.Error())           
        }

		var dto = map[uint]bool {pump.ID: gpio.AverageStateCache[pump.SensorPin]}
		return c.Status(200).JSON(dto)
    })

	// Check if Recipe can be served
	app.Get("/state/recipe/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		var recipe models.Recipe
		if tx := models.DB.First(&recipe, id); tx.Error != nil {
            return c.Status(404).SendString(tx.Error.Error())           
        }

		var dto = map[uint]bool {recipe.ID: gpio.CanBeServed(recipe)}
		return c.Status(200).JSON(dto)
    })

	// Same Check for all Recipes
	app.Get("/state/recipe", func(c *fiber.Ctx) error {
		var recipes []models.Recipe
		models.DB.Find(&recipes)

		var dto = make(map[uint]bool)

		for _, recipe := range recipes {
			dto[recipe.ID] = gpio.CanBeServed(recipe)
		}
		return c.Status(200).JSON(dto)
	})

	// Get the state of the current drink
    app.Get("/state/drink", func(c *fiber.Ctx) error {
		//TODO 
		drinkStatus := logic.DrinkStatus{ProgressInPercent: 40, Served: false, Canceled: false, Interrupted:  false}
		return c.Status(200).JSON(drinkStatus)
	})
}