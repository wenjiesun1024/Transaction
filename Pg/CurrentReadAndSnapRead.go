package pg

import (
	"Transaction/common"
	"Transaction/model"
	"sync"

	"gorm.io/gorm"
)

func PGCurrentReadAndSnapRead() {
	gormDB := common.InitPG(false)

	wg := sync.WaitGroup{}
	wg.Add(2)

	cond := &common.MyCond{
		Key:  int32(1),
		Cond: sync.NewCond(new(sync.Mutex)),
	}

	go func() {
		defer wg.Done()

		tx := gormDB.Begin() // Transaction 1
		defer tx.Commit()

		common.PrintlnAllData(tx, "1")
		common.WaitFor(cond, 1)

		common.WaitFor(cond, 4)

		//common.PrintlnAllData(tx, "3", clause.Locking{Strength: "UPDATE"}) // abort
		// common.PrintlnAllData(tx, "3", clause.Locking{Strength: "SHARE"}) // abort
		common.PrintlnAllData(tx, "3") // id=5, c=5, d=5

		// 自动检测更新丢失
		tx.Model(&model.T{}).Where("id = ?", 10).UpdateColumn("d", gorm.Expr("d + ?", 10)) // update id = 10, ok
		tx.Model(&model.T{}).Where("id = ?", 6).UpdateColumn("d", gorm.Expr("d + ?", 10))  // update id = 6, ok, no data, no update
		tx.Model(&model.T{}).Where("id = ?", 5).UpdateColumn("d", gorm.Expr("d + ?", 10))  // update id = 5, abort

		common.PrintlnAllData(tx, "4")
	}()

	go func() {
		defer wg.Done()
		tx := gormDB.Begin() // Transaction 2
		defer tx.Commit()

		common.WaitFor(cond, 2)

		tx.Model(&model.T{}).Where("id = ?", 5).Update("d", 100000)
		tx.Model(&model.T{}).Create(&model.T{ID: 6, C: 6, D: 6, E: 6})
		common.PrintlnAllData(tx, "2")

		common.WaitFor(cond, 3)
	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}
