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

	// Seed the Recipes
	// recipes := []Recipe {
	// 	{
	// 		Name: "Wodka-O",
	// 		Ingredients: []Ingredient {
	// 			{
	// 				Name: "Wodka",
	// 			},
	// 			{
	// 				Name: "Orangensaft",
	// 			},
	// 		},
	// 	},
	// }

	// // Create if not found
	// for _, recipe := range recipes {
	// 	var existingRecipe Recipe
	// 	res := db.Where(&Recipe{Name: recipe.Name}).First(&existingRecipe)

	// 	if res.RowsAffected == 0 {
	// 		db.Create(&recipe)
	// 	}
	// }
}

/*
// Create
db.Create(&Recipe{Code: "D42", Price: 100})

// Read
var recipe Recipe
db.First(&recipe, 1) // find recipe with integer primary key
db.First(&recipe, "code = ?", "D42") // find recipe with code D42

// Update - update recipe's price to 200
db.Model(&recipe).Update("Price", 200)
// Update - update multiple fields
db.Model(&recipe).Updates(Recipe{Price: 200, Code: "F42"}) // non-zero fields
db.Model(&recipe).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

// Delete - delete recipe
db.Delete(&recipe, 1)
}
*/
