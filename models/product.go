package models

import "gorm.io/gorm"

//Product -  model
type Product struct {
	gorm.Model
	Name      string
	Stock     int64 
	Price     float64
	Image     string 
}