package models

import (
	"gorm.io/gorm"
  )
  
  type Ingredient struct {
	gorm.Model
	Pump		Pump
	PumpID		int
	RecipeID	int
	Parts 		uint16	`gorm:"notNull;"`
  }

  func initIngredient() {
	// Migrate the schema
	DB.AutoMigrate(&Ingredient{})

	// Seed the Ingredients
	// ingredients := []Ingredient {
	// 	{
	// 		Name: "Orangensaft",
	// 		Pin: "D4",
	// 	},
	// 	{
	// 		Name: "Wodka",
	// 		Pin: "A3",
	// 	},
	// }

	// for _, ingredient := range ingredients {
	// 	var ingredientFoundOrCreated Ingredient
	// 	db.Where(ingredient).FirstOrCreate(&ingredientFoundOrCreated)
	// }
  }