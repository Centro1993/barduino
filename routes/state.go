package routes


import (
    "github.com/gofiber/fiber/v2"
    "barduino/models"
	"barduino/gpio"
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

		var dto = map[uint]bool {pump.SensorPin: gpio.AverageStateCache[pump.SensorPin]}
		return c.Status(200).JSON(dto)
    })

	// Get the state of the current drink
    app.Get("/state/drink", func(c *fiber.Ctx) error {
		// TODO
	})
}