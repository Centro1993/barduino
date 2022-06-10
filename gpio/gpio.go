package gpio

import (
	"barduino/models"
	"errors"
	"time"

	"github.com/stianeikeland/go-rpio/v4"
)

// TODO check if there are concurrency problems when using the Pins at the same time


const TIME_BETWEEN_SENSOR_CHECKS_IN_MS int64 = 100
const AMOUNT_OF_SENSOR_STATES_SAVED int = 10

// Access this to get the rolling average of the last few states of a Sensor
// key is the pump id for both maps
var AverageStateCache map[uint]bool
// This is used to save the last n sensor states internally
var sensorStateHistory map[uint][]rpio.State

func InitGPIO() error {
	err := rpio.Open()
	if err != nil {
		return errors.New("failed to access GPIO Pins")
	}

	go checkSensors()

	return nil
}

// This Function periodically gets all active Sensors, checks their state, and caches the last few states.
func checkSensors() {
	// init cache maps
	sensorStateHistory = make(map[uint][]rpio.State)
	AverageStateCache = make(map[uint]bool)

	for {
		// get all Sensor Pins
		var pumps []models.Pump
		models.DB.Find(&pumps)

		for _, pump := range pumps {
			// create stateHistory and pull down the resistor if the sensor has not been used yet
			if len(sensorStateHistory[pump.SensorPin]) == 0 {
				pin := rpio.Pin(pump.SensorPin)
				pin.PullDown()
				sensorStateHistory[pump.SensorPin] = make([]rpio.State, AMOUNT_OF_SENSOR_STATES_SAVED)
			}
			// read and append new state
			sensorStateHistory[pump.SensorPin] = append(sensorStateHistory[pump.SensorPin], readSensorState(pump.SensorPin))
			// Shift every other state one to the left, thereby deleting the oldest state
			sensorStateHistory[pump.SensorPin] = sensorStateHistory[pump.SensorPin][1:]
			// Compute Average State and cache it
			if state := computeAverageSensorState(&pump); state == rpio.High {
				AverageStateCache[pump.SensorPin] = true
			} else {
				AverageStateCache[pump.SensorPin] = false
			}
			
		}

		// Remove caches with no corresponding Pump (in case of Pump Deletion)
		for key, _ := range sensorStateHistory {
			pumpFound := false
			for _, pump := range pumps {
				if key == pump.SensorPin {
					pumpFound = true
					break
				}
			}
			if !pumpFound {
				delete(sensorStateHistory, key)
				delete(AverageStateCache, key)
				// reset sensor pin resistor
				pin := rpio.Pin(key)
				pin.PullOff()
			}
		}
		// check sleep repeat
		time.Sleep(time.Duration(TIME_BETWEEN_SENSOR_CHECKS_IN_MS * time.Hour.Milliseconds()))
	}
}

func readSensorState (sensorPin uint) rpio.State {
	pin := rpio.Pin(sensorPin)
	pin.Input()        // Input mode
	return pin.Read()  // Read state from pin (High / Low)
}

func computeAverageSensorState(pump *models.Pump) rpio.State {
	lastSensorStates := sensorStateHistory[pump.SensorPin]

	var highStateCount float64 = 0
	for _, state := range lastSensorStates {
		if state == rpio.High {
			highStateCount++
		}
	}

	if averageState := highStateCount / float64(len(lastSensorStates)); averageState >= 0.5 {
		return rpio.High
	} else {
		return rpio.Low
	}
}

func RunPump(barkeeper chan models.PumpStatus, pumpInstruction models.PumpInstruction) {
	pin := rpio.Pin(pumpInstruction.Pump.MotorPin)
	pin.Output()
	
	lastPumpStartTime := time.Now().UnixMilli()
	
	/* 	This is a BIDIRECTIONAL communication Loop
		Therefore, we have to send a message before we await an answer
		Else, we will enter a deadlock
	*/
	for pumpStatus := range barkeeper { 
		
		// cancel the drink if the barkeeper demands it
		if !pumpStatus.CurrentlyServing {
			pin.Low()
			barkeeper <- models.PumpStatus{
				CurrentlyServing: false,
				IngredientEmpty: false,
			}
			close(barkeeper)
			return
		}

		// pause the execution if another pump ran dry
		if pumpStatus.IngredientEmpty {
			pin.Low()
			// compute remaining time
			currentTime := time.Now().UnixMilli()
			pumpInstruction.TimeInMs -= (currentTime - lastPumpStartTime)
			// and assume the pump starts after this loop, so set the lastStartTime for the next remaining time computation
			lastPumpStartTime = currentTime
			// wait a while and check in with the barkeeper again
			time.Sleep(time.Duration(TIME_BETWEEN_SENSOR_CHECKS_IN_MS * int64(time.Millisecond)))
			barkeeper <- models.PumpStatus{
				CurrentlyServing: false,
				IngredientEmpty: false,
			}
			continue
		}

		// compute remaining time
		currentTime := time.Now().UnixMilli()
		pumpInstruction.TimeInMs -= (currentTime - lastPumpStartTime)

		// stop the execution if the pumping is done
		if pumpInstruction.TimeInMs <= 0 {
			break
		}

		// start the pump / keep it running
		lastPumpStartTime = time.Now().UnixMilli()
		pin.High()

		// If pumping the rest takes longer than until the next sensor check, wait until the next check
		if pumpInstruction.TimeInMs > TIME_BETWEEN_SENSOR_CHECKS_IN_MS {
			time.Sleep(time.Duration(TIME_BETWEEN_SENSOR_CHECKS_IN_MS * int64(time.Millisecond)))
			// and perform the sensor check
			if !AverageStateCache[pumpInstruction.Pump.SensorPin] {
				// if the sensor is low, stop running the pump
				pin.Low()

				// compute remaining time
				currentTime := time.Now().UnixMilli()
				pumpInstruction.TimeInMs -= (currentTime - lastPumpStartTime)

				// inform the barkeeper
				barkeeper <- models.PumpStatus{
					CurrentlyServing: true,
					IngredientEmpty: true,
				}

				// and continously check the sensors 
				for pumpStatus := range barkeeper {
					// stop if the barkeeper tells the pump to stop
					if !pumpStatus.CurrentlyServing {
						close(barkeeper)
						return
					}
					// else if the ingredient has been refilled
					if AverageStateCache[pumpInstruction.Pump.SensorPin] {
						// the pump starts up again, so set the lastStartTime for the next remaining time computation
						lastPumpStartTime = time.Now().UnixMilli()
						// and wait for the barkeeper to tell the pump to start up again
						break
					}
					// tell the barkeeper that we are still missing our ingredient
					barkeeper <- models.PumpStatus{
						CurrentlyServing: false,
						IngredientEmpty: true,
					}
					// wait until the next sensor check
					time.Sleep(time.Duration(TIME_BETWEEN_SENSOR_CHECKS_IN_MS * int64(time.Millisecond)))
				}
			}
		} else {	// else, sleep until the the Pumping is done, whatever is closer
			time.Sleep(time.Duration(pumpInstruction.TimeInMs * int64(time.Millisecond)))
		}
		// check in with the barkeeper and start the next loop
		barkeeper <- models.PumpStatus{
			CurrentlyServing: true,
			IngredientEmpty: false,
		}
	}

	pin.Low()

	barkeeper <- models.PumpStatus{
		CurrentlyServing: false,
		IngredientEmpty: false,
	}

	close(barkeeper)
}

func CanBeServed (recipe models.Recipe) bool {
	for _, ingredient := range recipe.Ingredients {
		if !AverageStateCache[ingredient.Pump.SensorPin] {
			return false
		}
	}
	return true
}