package mysql

import (
	"SQLIsolationLevelTest/common"
	"SQLIsolationLevelTest/model"
	"sync"
	"time"

	"gorm.io/gorm/clause"
)

func MysqlDeadLock() {
	gormDB := common.InitMysql()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()

		tx.Model(&model.User{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id", 3).First(&model.User{})
		time.Sleep(3 * time.Second)
		tx.Create(&model.User{Name: "Tom", Age: 44, UserID: 3})
	}()

	go func() {
		time.Sleep(2 * time.Second)

		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()

		common.PrintlnAllUsers(tx, "2", clause.Locking{Strength: "UPDATE"})

		// 在 RR 情况下， 因为 user_id 上有索引，所以这里只会锁住 user_id = 3 的行
		// 但是由于间隙锁的原因，会有(2, +∞)的间隙锁
		// 注意间隙锁之间是相互不冲突的， 与它冲突的是 “往这个间隙里插入一个新行” 的操作
		// 所以另外一个transaction 会被阻塞。同理，这个transaction 也会被阻塞直到另外一个transaction commit
		// 这就导致了死锁
		tx.Model(&model.User{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id", 3).First(&model.User{})
		tx.Create(&model.User{Name: "Tom", Age: 44, UserID: 3})
	}()

	wg.Wait()

	// 其中一个transaction 会被 deadlock rollback，所以还是会有一个transaction commit， 所以这里会打印出来Tom
	common.PrintlnAllUsers(gormDB, "end")
}
