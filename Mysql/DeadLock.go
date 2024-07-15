package mysql

import (
	"Transaction/common"
	"Transaction/model"
	"sync"

	"gorm.io/gorm/clause"
)

func MysqlDeadLock() {
	gormDB := common.InitMysql()

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

		tx.Model(&model.T{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id", 3).First(&model.T{})

		common.WaitFor(cond, 4)

		tx.Create(&model.T{ID: 3, C: 3, D: 3, E: 3})
	}()

	go func() {
		defer wg.Done()
		tx := gormDB.Begin() // Transaction 2
		defer tx.Commit()

		common.WaitFor(cond, 2)

		common.PrintlnAllData(tx, "2")

		// 在 RR 情况下， 因为 user_id 上有索引，由于间隙锁的原因，会有(1, 5)的间隙锁
		// 注意间隙锁之间是相互不冲突的， 与它冲突的是 “往这个间隙里插入一个新行” 的操作
		tx.Model(&model.T{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id", 4).First(&model.T{})

		common.WaitFor(cond, 3)

		// 所以另外一个transaction 会被阻塞。同理，这个transaction 也会被阻塞直到另外一个transaction commit
		// 这就导致了死锁
		tx.Create(&model.T{ID: 4, C: 4, D: 4, E: 4})
	}()

	wg.Wait()

	// 其中一个transaction 会被 deadlock rollback，所以还是会有一个transaction commit
	common.PrintlnAllData(gormDB, "end")
}
