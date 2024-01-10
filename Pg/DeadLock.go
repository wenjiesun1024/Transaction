package pg

import (
	"Transaction/common"
	"Transaction/model"
	"database/sql"
	"sync"
	"time"

	"gorm.io/gorm/clause"
)

func PGDeadLock() {
	gormDB := common.InitMysql()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		tx := gormDB.Begin(&sql.TxOptions{
			Isolation: sql.LevelRepeatableRead,
		})
		defer tx.Commit()

		common.PrintlnAllData(tx, "1")

		tx.Model(&model.T{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id", 3).First(&model.T{})
		time.Sleep(3 * time.Second)
		tx.Create(&model.T{ID: 3, C: 3, D: 3})
	}()

	go func() {
		defer wg.Done()

		time.Sleep(2 * time.Second)

		tx := gormDB.Begin(&sql.TxOptions{
			Isolation: sql.LevelRepeatableRead,
		})
		defer tx.Commit()

		common.PrintlnAllData(tx, "2")

		// 在 RR 情况下， 因为 user_id 上有索引，由于间隙锁的原因，会有(2, +∞)的间隙锁
		// 注意间隙锁之间是相互不冲突的， 与它冲突的是 “往这个间隙里插入一个新行” 的操作
		tx.Model(&model.T{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id", 4).First(&model.T{})

		// 所以另外一个transaction 会被阻塞。同理，这个transaction 也会被阻塞直到另外一个transaction commit
		// 这就导致了死锁
		tx.Create(&model.T{ID: 4, C: 4, D: 4})
	}()

	wg.Wait()

	// 其中一个transaction 会被 deadlock rollback，所以还是会有一个transaction commit， 所以这里会打印出来Tom
	common.PrintlnAllData(gormDB, "end")
}
