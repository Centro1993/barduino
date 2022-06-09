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
	pin.High()
	
	//TODO communicate with barkeeper before each cycle
	//TODO take delays into account
	//TODO return progress
	// Sleep until the next Sensor Check is due or until finished
	var timeEclipsed int64 = 0
	startTime := time.Now().UnixMilli()
	for pumpInstruction.TimeInMs > timeEclipsed {
		timeRemaining := pumpInstruction.TimeInMs - timeEclipsed
		// Sleep until the next Sensor Check is due
		if timeRemaining > TIME_BETWEEN_SENSOR_CHECKS_IN_MS {
			time.Sleep(time.Duration(TIME_BETWEEN_SENSOR_CHECKS_IN_MS * int64(time.Millisecond)))
			// and either perform the sensor check
			if sensorState := AverageStateCache[pumpInstruction.Pump.SensorPin]; !sensorState{
				// if the sensor is low, stop running the pump
				pin.Low()
				// inform the barkeeper
				barkeeper <- models.PumpStatus{
					CurrentlyServing: true,
					IngredientEmpty: true,
				}
				// and wait until the barkeeper says to continue or stop
				for pumpStatus := range barkeeper {
					if !pumpStatus.CurrentlyServing {
						close(barkeeper)
						return
					}
					if !pumpStatus.IngredientEmpty {
						pin.High()
						break
					}
				}
			}
		} else {	// or sleep until the the Pumping is done, whatever is closer
			time.Sleep(time.Duration(timeRemaining * int64(time.Millisecond)))
		}
		// compute remaining time
		timeEclipsed = pumpInstruction.TimeInMs + startTime - time.Now().UnixMilli()
	}

	pin.Low()

	barkeeper <- models.PumpStatus{
		ProgressInPercent: 100,
		CurrentlyServing: false,
		IngredientEmpty: false,
	}

	close(barkeeper)
}

func StopPump(pump models.Pump) {
	pin := rpio.Pin(pump.MotorPin)
	pin.Low()
}

func CanBeServed (recipe models.Recipe) bool {
	for _, ingredient := range recipe.Ingredients {
		if !AverageStateCache[ingredient.Pump.SensorPin] {
			return false
		}
	}
	return true
}