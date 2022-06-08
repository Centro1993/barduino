package models

import (
	"gorm.io/gorm"
)

type Drink struct {
	gorm.Model
	Recipe		Recipe	`gorm:"notNull"`
	Served	bool		`gorm:"default:false"`
}

func initDrink() {
	// Migrate the schema
	DB.AutoMigrate(&Drink{})
}