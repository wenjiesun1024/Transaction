package pg

import (
	"Transaction/common"
	"Transaction/model"
	"sync"
	"time"

	"gorm.io/gorm/clause"
)

func PGDeadLock() {
	gormDB := common.InitPG(false)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		tx := gormDB.Begin()
		defer tx.Commit()

		common.PrintlnAllData(tx, "1")

		tx.Model(&model.T{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id", 3).First(&model.T{})
		time.Sleep(3 * time.Second)
		tx.Create(&model.T{ID: 3, C: 3, D: 3, E: 3})
	}()

	go func() {
		defer wg.Done()

		time.Sleep(2 * time.Second)

		tx := gormDB.Begin()
		defer tx.Commit()

		common.PrintlnAllData(tx, "2")

		tx.Model(&model.T{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id", 4).First(&model.T{})

		tx.Create(&model.T{ID: 4, C: 4, D: 4, E: 4})
	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}
