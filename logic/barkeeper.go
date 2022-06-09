package logic

import (
	"barduino/gpio"
	"barduino/models"
	"errors"
)

const SERVING_SIZE_IN_ML uint = 300

var drinkStatus models.DrinkStatus = models.DrinkStatus{
	CurrentlyServing: false,
}
var runningPumps map[uint]chan models.PumpStatus

func PourDrink (drink models.Drink) error {
	recipe := drink.Recipe

	if drinkStatus.CurrentlyServing {
		return errors.New("cannot pour drink, currently serving another drink")
	}
	if !gpio.CanBeServed(recipe) {
		return errors.New("cannot pour drink, ingredient missing")
	}

	drinkStatus.CurrentlyServing = true
	pumpInstructions := recipe.ConvertRecipeToPumpInstructions(SERVING_SIZE_IN_ML)
	for _, pumpInstruction := range pumpInstructions {
		runningPumps[pumpInstruction.Pump.ID] = make(chan models.PumpStatus)
		go gpio.RunPump(runningPumps[pumpInstruction.Pump.ID], pumpInstruction)
		go monitorPump(pumpInstruction.Pump)
	}

	return nil
}

func CancelDrink () error {
	if !drinkStatus.CurrentlyServing {
		return errors.New("cannot cancel drink, not currently serving")
	}

	drinkStatus.CurrentlyServing = false
	return nil
}

//TODO monitor Progress
func monitorPump(pump models.Pump) error {
	if runningPumps[pump.ID] == nil {
		return errors.New("pump has no open channel and cannot be monitored")
	}

	for pumpStatus := range runningPumps[pump.ID] {
		// this stops all pumps if this pump runs dry
		if pumpStatus.IngredientEmpty {
			drinkStatus.IngredientEmpty = true
		}
		// return the overall drink status to the pump
		runningPumps[pump.ID] <- models.PumpStatus{
			IngredientEmpty: drinkStatus.IngredientEmpty,
			CurrentlyServing: drinkStatus.CurrentlyServing,
		}
	}

	return nil
}

func GetDrinkStatus() models.DrinkStatus {
	return drinkStatus
}