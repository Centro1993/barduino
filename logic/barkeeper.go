package logic

import (
	"barduino/gpio"
	"barduino/models"
	"errors"
	"fmt"
)

const SERVING_SIZE_IN_ML uint = 200

// Drink serving state
var drinkStatus models.DrinkStatus = models.DrinkStatus{
	CurrentlyServing: false,
	IngredientEmpty: false,
}
// Pump channels
var runningPumps map[uint]chan models.PumpStatus

func PourDrink (drink models.Drink) error {
	runningPumps = make(map[uint]chan models.PumpStatus)
	recipe := drink.Recipe

	if drinkStatus.CurrentlyServing {
		return errors.New("cannot pour drink, currently serving another drink")
	}
	if !gpio.CanBeServed(recipe) {
		return errors.New("cannot pour drink, ingredient missing")
	}

	drinkStatus.CurrentlyServing = true
	drinkStatus.IngredientEmpty = false 

	pumpInstructions := recipe.ConvertRecipeToPumpInstructions(SERVING_SIZE_IN_ML)
	for _, pumpInstruction := range pumpInstructions {
		runningPumps[pumpInstruction.Pump.ID] = make(chan models.PumpStatus)
		go gpio.RunPump(runningPumps[pumpInstruction.Pump.ID], pumpInstruction)
		go monitorPump(pumpInstruction.Pump)
	}

	return nil
}

func CancelDrink () error {
	fmt.Println("cancel drink")
	fmt.Print(drinkStatus)
	drinkStatus.CurrentlyServing = false

	// As a safety Measure, manually disable all pumps
	var pumps []models.Pump
	models.DB.Find(&pumps)

	for _, pump := range pumps {
		gpio.StopMotor(pump)
	}

	if !drinkStatus.CurrentlyServing {
		return errors.New("cannot cancel drink, not currently serving")
	}
	return nil
}

func monitorPump(pump models.Pump) error {
	if runningPumps[pump.ID] == nil {
		return errors.New("pump has no open channel and cannot be monitored")
	}
	currentChannel := runningPumps[pump.ID]

	// tell pump to start
	currentChannel <- models.PumpStatus{
		CurrentlyServing: drinkStatus.CurrentlyServing,
		IngredientEmpty: drinkStatus.IngredientEmpty,
	}

	/* 	This is a BIDIRECTIONAL communication Loop
		Therefore, we have to send a message before we await an answer
		Else, we will enter a deadlock
	*/
	for pumpStatus := range currentChannel {
		// this stops all pumps if this pump runs dry
		if pumpStatus.IngredientEmpty {
			drinkStatus.IngredientEmpty = true
			// acknowledge the report
			currentChannel <- models.PumpStatus{
				IngredientEmpty: drinkStatus.IngredientEmpty,
				CurrentlyServing: drinkStatus.CurrentlyServing,
			}
			// continously ask the pump if the ingredient has been refilled
			for dryPumpStatus := range currentChannel {
				// if the pump has been refilled, inform the other pumpMonitors
				if !dryPumpStatus.IngredientEmpty {
					drinkStatus.IngredientEmpty = false
					// and continue pumping
					break
				}
				// if not, tell the pump to keep checking (or to stop if the user chose so)
				currentChannel <- models.PumpStatus{
					IngredientEmpty: drinkStatus.IngredientEmpty,
					CurrentlyServing: drinkStatus.CurrentlyServing,
				}
			}
		}
		// return the overall drink status to the pump
		currentChannel <- models.PumpStatus{
			IngredientEmpty: drinkStatus.IngredientEmpty,
			CurrentlyServing: drinkStatus.CurrentlyServing,
		}
	}
	// When the channel has been closed, delete it from the channel map
	delete(runningPumps, pump.ID)

	// if there are no more channels in the map, we are no longer serving
	if (len(runningPumps) == 0) {
		drinkStatus.CurrentlyServing = false
	}

	return nil
}

func GetDrinkStatus() models.DrinkStatus {
	return drinkStatus
}