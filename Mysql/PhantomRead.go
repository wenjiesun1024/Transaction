package mysql

import (
	"Transaction/common"
	"Transaction/model"
	"sync"
)

func MysqlPhantomRead() {
	gormDB := common.InitMysql()

	wg := sync.WaitGroup{}
	wg.Add(2)

	cond := &common.MyCond{
		Key:  int32(1),
		Cond: sync.NewCond(new(sync.Mutex)),
	}

	go func() {
		defer wg.Done()

		tx := gormDB.Begin()
		defer tx.Commit()

		common.PrintlnAllData(tx, "1")

		common.WaitFor(cond, 1)

		common.WaitFor(cond, 4)

		tx.Model(&model.T{}).Where("1=1").Update("d", 100000)

		common.PrintlnAllData(tx, "2")
	}()

	go func() {
		defer wg.Done()

		common.WaitFor(cond, 2)

		gormDB.Model(&model.T{}).Create(&model.T{ID: 6, C: 6, D: 6})

		common.WaitFor(cond, 3)

	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}
