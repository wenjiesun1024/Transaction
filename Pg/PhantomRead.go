package pg

import (
	"Transaction/common"
	"Transaction/model"
	"database/sql"

	// "database/sql"
	"sync"
	"time"
)

func PGPhantomRead() {
	gormDB := common.InitPG()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		tx := gormDB.Begin(&sql.TxOptions{
			Isolation: sql.LevelRepeatableRead,
		})
		defer tx.Commit()

		common.PrintlnAllData(tx, "1")

		time.Sleep(3 * time.Second)

		tx.Model(&model.T{}).Where("1=1").Update("d", 100000)

		common.PrintlnAllData(tx, "2")
	}()

	go func() {
		defer wg.Done()

		time.Sleep(2 * time.Second)

		gormDB.Model(&model.T{}).Create(&model.T{ID: 6, C: 6, D: 6})
	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}
