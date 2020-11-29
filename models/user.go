package models

import (
	"gorm.io/gorm"
)

//User - model
type User struct {
	gorm.Model
	Username string `gorm:"unique" form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
	Level    string `gorm:"Default:normal"`
}
