package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	UserID int `gorm:"uniqueIndex;column:user_id"`
	Name   string
	Age    int
}
