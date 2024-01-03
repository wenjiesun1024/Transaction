package model

type T struct {
	ID int `gorm:"primary_key"`
	C  int `gorm:"index"`
	D  int
}
