package pg

import (
	"Transaction/common"
	"Transaction/model"
	"sync"
	"time"
)

// 自动检测更新
func PGUpdate() {
	gormDB := common.InitPG()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		tx := gormDB.Begin()
		defer tx.Commit()

		common.PrintlnAllData(tx, "1") // read old value

		time.Sleep(3 * time.Second)

		tx.Model(&model.T{}).Where("id = ?", 10).Update("d", 500) // update, abort

		common.PrintlnAllData(tx, "3")
	}()

	go func() {
		defer wg.Done()
		time.Sleep(2 * time.Second)

		gormDB.Model(&model.T{}).Where("id = ?", 10).Update("d", 500)
		common.PrintlnAllData(gormDB, "2")
	}()

	wg.Wait()
	common.PrintlnAllData(gormDB, "end")
}

// type A struct {
// 	ID int `gorm:"primaryKey"`
// }

// func PGUpdate2() {
// 	gormDB := common.InitPG()

// 	gormDB.AutoMigrate(&A{})

// 	wg := sync.WaitGroup{}
// 	wg.Add(2)

// 	go func() {
// 		defer wg.Done()

// 		tx := gormDB.Begin()
// 		defer tx.Commit()

// 		tx.Model(&A{}).Where("id = ?", 1000000).First(&A{})

// 		time.Sleep(3 * time.Second)

// 		tx.Model(&model.T{}).Where("id = ?", 10).Update("d", 1000) // abort

// 		common.PrintlnAllData(tx, "3")
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		time.Sleep(2 * time.Second)

// 		gormDB.Model(&model.T{}).Where("id = ?", 10).Update("d", 500)
// 		common.PrintlnAllData(gormDB, "2")
// 	}()

// 	wg.Wait()
// 	common.PrintlnAllData(gormDB, "end")
// }
