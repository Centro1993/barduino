package models

import (
	"gorm.io/gorm"
  )
  
  type Pump struct {
	gorm.Model
	Name string	`gorm:"unique;notNull;"`
	// TODO unique over two rows
	MotorPin string	`gorm:"index:,unique,composite:pin;notNull;"`
	SensorPin string `gorm:"index:,unique,composite:pin;notNull;"`
	MlPerMinute uint `gorm:"default:48"`
  }

  func initPump() {
	// Migrate the schema
	DB.AutoMigrate(&Pump{})

	// Seed the Pumps
	// pumps := []Pump {
	// 	{
	// 		Name: "Orangensaft",
	// 		Pin: "D4",
	// 	},
	// 	{
	// 		Name: "Wodka",
	// 		Pin: "A3",
	// 	},
	// }

	// for _, pump := range pumps {
	// 	var pumpFoundOrCreated Pump
	// 	db.Where(pump).FirstOrCreate(&pumpFoundOrCreated)
	// }
  }