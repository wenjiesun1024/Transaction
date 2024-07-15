package common

import (
	"fmt"
	"sort"
	"sync"
	"sync/atomic"

	"Transaction/model"

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

	// drop table
	gormDB.Migrator().DropTable(&model.T{})

	// use gorm to create a table
	gormDB.AutoMigrate(&model.T{})

	// insert data
	gormDB.Create(&model.T{ID: 1, C: 1, D: 1, E: 1})
	gormDB.Create(&model.T{ID: 5, C: 5, D: 5, E: 5})
	gormDB.Create(&model.T{ID: 10, C: 10, D: 10, E: 10})
	gormDB.Create(&model.T{ID: 15, C: 15, D: 15, E: 15})
	gormDB.Create(&model.T{ID: 20, C: 20, D: 20, E: 20})
	gormDB.Create(&model.T{ID: 25, C: 25, D: 25, E: 25})

	PrintlnAllData(gormDB, "Init")

	return gormDB
}

func InitPG(EnableRR bool) *gorm.DB {
	// use gorm to creat a table
	gormDB, err := gorm.Open(postgres.Open("host=localhost user=root password=pass dbname=db port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// set transaction isolation level
	if EnableRR {
		err = gormDB.Exec("SET SESSION CHARACTERISTICS AS TRANSACTION ISOLATION LEVEL REPEATABLE READ").Error
		if err != nil {
			panic(err)
		}
	}

	// drop table
	gormDB.Migrator().DropTable(&model.T{})

	// use gorm to create a table
	gormDB.AutoMigrate(&model.T{})

	// insert data
	gormDB.Create(&model.T{ID: 1, C: 1, D: 1, E: 1})
	gormDB.Create(&model.T{ID: 5, C: 5, D: 5, E: 5})
	gormDB.Create(&model.T{ID: 10, C: 10, D: 10, E: 10})
	gormDB.Create(&model.T{ID: 15, C: 15, D: 15, E: 15})
	gormDB.Create(&model.T{ID: 20, C: 20, D: 20, E: 20})
	gormDB.Create(&model.T{ID: 25, C: 25, D: 25, E: 25})

	PrintlnAllData(gormDB, "Init")
	return gormDB
}

func PrintlnAllData(db *gorm.DB, tag string, Clauses ...clause.Expression) error {
	var T []model.T
	if err := db.Clauses(Clauses...).Find(&T).Error; err != nil {
		return err
	}
	sort.Slice(T, func(i, j int) bool {
		if T[i].ID == T[j].ID {
			return T[i].C < T[j].C
		}
		return T[i].ID < T[j].ID
	})
	fmt.Println(tag)
	for _, i := range T {
		fmt.Printf("%+v\n", i)
	}
	fmt.Println("------------------------")
	return nil
}

type MyCond struct {
	*sync.Cond
	Key int32
}

func WaitFor(cond *MyCond, expectedValue int32) {
	cond.L.Lock()
	for cond.Key != expectedValue {
		cond.Wait()
	}
	atomic.AddInt32(&cond.Key, 1)
	cond.L.Unlock()
	cond.Broadcast()
}
