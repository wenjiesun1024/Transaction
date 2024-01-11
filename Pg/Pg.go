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

		common.PrintlnAllData(tx, "1") // read old value(remove this, update success)

		time.Sleep(3 * time.Second)

		tx.Model(&model.T{}).Where("id = ?", 10).Update("d", 500) // update, abort(if we read old value, else update success)

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
