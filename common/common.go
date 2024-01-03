package common

import (
	"fmt"

	"SQLIsolationLevelTest/model"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InitMysql() *gorm.DB {
	// use gorm to connect mysql
	dsn := "root:pass@tcp(localhost:3306)/db?charset=utf8mb4&parseTime=True&loc=Local"
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// use gorm to create a table
	gormDB.AutoMigrate(&model.T{})

	// clear all data
	gormDB.Unscoped().Where("1 = 1").Delete(&model.T{})

	// insert data
	gormDB.Create(&model.T{ID: 1, C: 1, D: 1})
	gormDB.Create(&model.T{ID: 5, C: 5, D: 5})
	gormDB.Create(&model.T{ID: 10, C: 10, D: 10})
	gormDB.Create(&model.T{ID: 15, C: 15, D: 15})
	gormDB.Create(&model.T{ID: 20, C: 20, D: 20})
	gormDB.Create(&model.T{ID: 25, C: 25, D: 25})

	return gormDB
}

func InitPG() *gorm.DB {
	// use gorm to creat a table
	gormDB, err := gorm.Open(postgres.Open("host=localhost user=root password=pass dbname=db port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// set transaction isolation level
	err = gormDB.Exec(`set transaction isolation level repeatable read`).Error
	if err != nil {
		panic(err)
	}

	// use gorm to create a table
	gormDB.AutoMigrate(&model.T{})

	// clear all data
	gormDB.Unscoped().Where("1 = 1").Delete(&model.T{})
	// insert data
	gormDB.Create(&model.T{ID: 1, C: 1, D: 1})
	gormDB.Create(&model.T{ID: 5, C: 5, D: 5})
	gormDB.Create(&model.T{ID: 10, C: 10, D: 10})
	gormDB.Create(&model.T{ID: 15, C: 15, D: 15})
	gormDB.Create(&model.T{ID: 20, C: 20, D: 20})
	gormDB.Create(&model.T{ID: 25, C: 25, D: 25})

	return gormDB
}

func PrintlnAllData(db *gorm.DB, tag string, Clauses ...clause.Expression) {
	var T []model.T
	db.Clauses(Clauses...).Find(&T)
	fmt.Println(tag)
	for _, i := range T {
		fmt.Println(i)
	}
	fmt.Println("------------------------")
}
