package gpio

import (
	"errors"
	"time"

	"barduino/models"

	"github.com/stianeikeland/go-rpio/v4"
)

// TODO check if there are concurrency problems when using the Pins at the same time


const TIME_BETWEEN_SENSOR_CHECKS_IN_MS int64 = 100
const AMOUNT_OF_SENSOR_STATES_SAVED int = 10

// Access this to get the rolling average of the last few states of a Sensor
var AverageStateCache map[uint]bool
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
			// create slice if it does not exist yet
			if len(sensorStateHistory[pump.SensorPin]) == 0 {
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

func RunPump(pump *models.Pump, timeInMs int64) error {
	pin := rpio.Pin(pump.MotorPin)
	pin.Output()
	pin.High()
	
	// Sleep until the next Sensor Check is due or until finished
	var timeEclipsed int64 = 0
	startTime := time.Now().UnixMilli()
	for timeInMs > timeEclipsed {
		timeRemaining := timeInMs - timeEclipsed
		// Sleep until the next Sensor Check is due
		if timeRemaining > TIME_BETWEEN_SENSOR_CHECKS_IN_MS {
			time.Sleep(time.Duration(TIME_BETWEEN_SENSOR_CHECKS_IN_MS * int64(time.Millisecond)))
			// and perform the sensor check
			if sensorState := AverageStateCache[pump.SensorPin]; !sensorState{
				// if the sensor is low, stop running the pump
				pin.Low()
				return errors.New("ingredient is empty")
			}
		} else {	// or sleep until the the Pumping is done, whatever is closer
			time.Sleep(time.Duration(timeRemaining * int64(time.Millisecond)))
		}
		// compute remaining time
		timeEclipsed = timeInMs + startTime - time.Now().UnixMilli()
	}

	pin.Low()

	return nil
}