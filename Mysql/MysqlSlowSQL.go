package mysql

import (
	"SQLIsolationLevelTest/common"
	"sync"
	"time"

	"gorm.io/gorm/clause"
)

func MysqlSlowSQL() {
	gormDB := common.InitMysql()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()

		// 在 RR，for update 会锁住所有读过的行，所以这里会锁住所有的行， 同时所有区间上的间隙也会被锁住
		// 另外一个transaction 已经锁住所有的行，所以 这个transaction 会被阻塞直到另外一个transaction commit
		common.PrintlnAllUsers(tx, "1", clause.Locking{Strength: "UPDATE"})
		time.Sleep(10 * time.Second)
	}()

	go func() {
		time.Sleep(2 * time.Second)

		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()

		common.PrintlnAllUsers(tx, "2", clause.Locking{Strength: "UPDATE"})
	}()

	wg.Wait()

	common.PrintlnAllUsers(gormDB, "end")
}
