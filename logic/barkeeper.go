package logic

import (
	"barduino/gpio"
	"barduino/models"
	"errors"
	"sync"
)

const SERVING_SIZE_IN_ML uint = 200

// Drink serving state
var drinkStatus models.DrinkStatus = models.DrinkStatus{
	CurrentlyServing: false,
	IngredientEmpty: false,
}
var drinkStatusMutex sync.RWMutex = sync.RWMutex{}

// Pump channels
var runningPumps map[uint]chan models.PumpStatus
var runningPumpsMutex sync.RWMutex = sync.RWMutex{}

func PourDrink (drink models.Drink) error {
	runningPumpsMutex.Lock()
	runningPumps = make(map[uint]chan models.PumpStatus)
	runningPumpsMutex.Unlock()
	recipe := drink.Recipe

	drinkStatusMutex.Lock()
	if drinkStatus.CurrentlyServing {
		return errors.New("cannot pour drink, currently serving another drink")
	}
	if !gpio.CanBeServed(recipe) {
		return errors.New("cannot pour drink, ingredient missing")
	}

	drinkStatus.CurrentlyServing = true
	drinkStatus.IngredientEmpty = false 
	drinkStatusMutex.Unlock()

	pumpInstructions := recipe.ConvertRecipeToPumpInstructions(SERVING_SIZE_IN_ML)
	for _, pumpInstruction := range pumpInstructions {
		runningPumps[pumpInstruction.Pump.ID] = make(chan models.PumpStatus)
		go gpio.RunPump(runningPumps[pumpInstruction.Pump.ID], pumpInstruction)
		go monitorPump(pumpInstruction.Pump)
	}

	return nil
}

func CancelDrink () error {
	drinkStatusMutex.Unlock()
	if !drinkStatus.CurrentlyServing {
		return errors.New("cannot cancel drink, not currently serving")
	}

	drinkStatus.CurrentlyServing = false
	drinkStatusMutex.Lock()

	// As a safety Measure, manually disable all pumps
	var pumps []models.Pump
	models.DB.Find(&pumps)

	for _, pump := range pumps {
		gpio.StopMotor(pump)
	}
	return nil
}

func monitorPump(pump models.Pump) error {
	runningPumpsMutex.Lock()
	if runningPumps[pump.ID] == nil {
		return errors.New("pump has no open channel and cannot be monitored")
	}
	currentChannel := runningPumps[pump.ID]
	runningPumpsMutex.Unlock()

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
			drinkStatusMutex.Lock()
			drinkStatus.IngredientEmpty = true
			// acknowledge the report
			currentChannel <- models.PumpStatus{
				IngredientEmpty: drinkStatus.IngredientEmpty,
				CurrentlyServing: drinkStatus.CurrentlyServing,
			}
			drinkStatusMutex.Unlock()

			// continously ask the pump if the ingredient has been refilled
			for dryPumpStatus := range currentChannel {
				// if the pump has been refilled, inform the other pumpMonitors
				drinkStatusMutex.Unlock()
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
				drinkStatusMutex.Lock()
			}
		}
		// return the overall drink status to the pump
		drinkStatusMutex.Unlock()
		currentChannel <- models.PumpStatus{
			IngredientEmpty: drinkStatus.IngredientEmpty,
			CurrentlyServing: drinkStatus.CurrentlyServing,
		}
		drinkStatusMutex.Lock()
	}
	// When the channel has been closed, delete it from the channel map
	runningPumpsMutex.Lock()
	delete(runningPumps, pump.ID)

	// if there are no more channels in the map, we are no longer serving
	if (len(runningPumps) == 0) {
		drinkStatus.CurrentlyServing = false
	}
	runningPumpsMutex.Unlock()

	return nil
}

func GetDrinkStatus() models.DrinkStatus {
	return drinkStatus
}