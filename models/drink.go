package models

import (
	"gorm.io/gorm"
)

type Drink struct {
	gorm.Model
	Recipe		Recipe	`gorm:"notNull"`
	RecipeID	int
}

type DrinkStatus struct{
	IngredientEmpty		bool
	CurrentlyServing	bool
}

func initDrink() {
	// Migrate the schema
	DB.AutoMigrate(&Drink{})
}