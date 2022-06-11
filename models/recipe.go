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

type PumpInstruction struct{
	Pump		Pump
	TimeInMs	int64
}

func (recipe Recipe) ConvertRecipeToPumpInstructions(servingSizeInMl uint) []PumpInstruction {
	var partsTotal uint = 0

	for _, ingredients := range recipe.Ingredients {
		partsTotal += ingredients.Parts
	}

	var mlPerPart uint = servingSizeInMl / partsTotal

	var pumpInstructions []PumpInstruction = make([]PumpInstruction, len(recipe.Ingredients))
	for i, ingredient := range recipe.Ingredients {
		pumpInstructions[i] = PumpInstruction{
			Pump: ingredient.Pump,
			TimeInMs: int64(float64(ingredient.Parts) * float64(mlPerPart) / (float64(ingredient.Pump.MlPerMinute) / 60.0 / 1000.0)),
		}
	}

	return pumpInstructions
}