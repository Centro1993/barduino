package models

import (
	"gorm.io/gorm"
)

type Recipe struct {
	gorm.Model
	Name        string		 `gorm:"unique;notNull;"`
	Ingredients []Ingredient
}

func initRecipe() {
	// Migrate the schema
	DB.AutoMigrate(&Recipe{})
}

type pumpInstruction struct{
	pump		Pump
	timeInMs	uint
}

func (Recipe) ConvertRecipeToPumpInstructions(recipe *Recipe, servingSizeInMl uint) []pumpInstruction {
	var partsTotal uint = 0

	for _, ingredients := range recipe.Ingredients {
		partsTotal += ingredients.Parts
	}

	var mlPerPart uint = servingSizeInMl / partsTotal

	var pumpInstructions []pumpInstruction = make([]pumpInstruction, len(recipe.Ingredients))
	for i, ingredient := range recipe.Ingredients {
		pumpInstructions[i] = pumpInstruction{
			pump: ingredient.Pump,
			timeInMs: uint(float64(ingredient.Parts) * float64(mlPerPart) / (float64(ingredient.Pump.MlPerMinute) / 60.0 / 10000.0)),
		}
	}

	return pumpInstructions
}