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
	ProgressInPercent	uint
	IngredientEmpty		bool
	CurrentlyServing	bool
	Canceled			bool
}

func initDrink() {
	// Migrate the schema
	DB.AutoMigrate(&Drink{})
}