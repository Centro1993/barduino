package models

import (
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
  )


var DB *gorm.DB

func InitDabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("barduino.sqlite3"), &gorm.Config{})
	if err != nil {
	  panic("failed to connect database")
	}

	// Order is important here, Relationships have to be set up correctly
	initPump()
	initIngredient()
	initRecipe()
}