package models

import (
	"gorm.io/gorm"
  )
  
  type Pump struct {
	gorm.Model
	Name string	`gorm:"unique;notNull;"`
	// TODO unique over two rows
	MotorPin uint	`gorm:"index:,unique,composite:pin;notNull;"`
	SensorPin uint `gorm:"index:,unique,composite:pin;notNull;"`
	MlPerMinute uint `gorm:"default:150"`
  }

  type PumpStatus struct {
	  CurrentlyServing	bool
	  IngredientEmpty 	bool
  }

  func initPump() {
	// Migrate the schema
	DB.AutoMigrate(&Pump{})
  }