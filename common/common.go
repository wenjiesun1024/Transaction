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
	gormDB.AutoMigrate(&model.User{})

	// clear all data
	gormDB.Unscoped().Where("1 = 1").Delete(&model.User{})

	// insert data
	gormDB.Create(&model.User{Name: "Bob", Age: 30, UserID: 1})
	gormDB.Create(&model.User{Name: "John", Age: 40, UserID: 2})

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
	gormDB.AutoMigrate(&model.User{})

	// clear all data
	gormDB.Unscoped().Where("1 = 1").Delete(&model.User{})
	// insert data
	gormDB.Create(&model.User{Name: "Bob", Age: 30, UserID: 1})
	gormDB.Create(&model.User{Name: "John", Age: 40, UserID: 2})

	return gormDB
}

func PrintlnAllUsers(db *gorm.DB, tag string, Clauses ...clause.Expression) {
	var users []model.User
	db.Clauses(Clauses...).Find(&users)
	fmt.Println(tag)
	for _, user := range users {
		fmt.Println(user.Name, user.Age, user.UserID)
	}
	fmt.Println("------------------------")
}
