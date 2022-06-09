package models

import (
	"gorm.io/gorm"
  )
  
  type Ingredient struct {
	gorm.Model
	Pump		Pump
	PumpID		int
	RecipeID	int
	Parts 		uint	`gorm:"notNull;"`
  }

  func initIngredient() {
	// Migrate the schema
	DB.AutoMigrate(&Ingredient{})
  }