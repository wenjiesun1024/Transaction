package mysql

import (
	"SQLIsolationLevelTest/common"
	"SQLIsolationLevelTest/model"
	"sync"
	"time"

	"gorm.io/gorm/clause"
)

func MysqlLock() {
	gormDB := common.InitMysql()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		tx := gormDB.Begin()
		defer tx.Commit()

		tx.Model(&model.T{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", 7).First(&model.T{}) // 间隙锁(5, 10)
		common.PrintlnAllData(tx, "1")

		time.Sleep(10 * time.Second)
	}()

	go func() {
		defer wg.Done()

		time.Sleep(2 * time.Second)
		gormDB.Create(&model.T{ID: 8, C: 8, D: 8}) // block
		common.PrintlnAllData(gormDB, "2")

	}()

	go func() {
		defer wg.Done()

		time.Sleep(3 * time.Second)
		gormDB.Model(&model.T{}).Where("id = ?", 10).Update("d", 888) // not block
		common.PrintlnAllData(gormDB, "3")
	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}

func MysqlLock2() {
	gormDB := common.InitMysql()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		tx := gormDB.Begin()
		defer tx.Commit()

		tx.Debug().Raw("select id from ts where c = ? lock in share mode", 5).Scan(&model.T{}) // 间隙锁(0, 10), 注意锁的是覆盖索引
		// tx.Debug().Raw("select id from ts where c = ? for update", 5).Scan(&model.T{}) // 间隙锁(0, 10), 注意会锁住主键索引。所以此时更新也会被阻塞
		// tx.Debug().Raw("select d from ts where c = ? lock in share mode", 5).Scan(&model.T{}) // 间隙锁(0, 10), 注意覆盖索引中没有 d，所以回表的过程中，主键索引也会被锁住
		common.PrintlnAllData(tx, "1")

		time.Sleep(6 * time.Second)
	}()

	go func() {
		defer wg.Done()

		time.Sleep(2 * time.Second)
		gormDB.Debug().Model(&model.T{}).Where("id = ?", 5).Update("d", 888) // not block, 因为锁的是覆盖索引
		common.PrintlnAllData(gormDB, "2")

	}()

	go func() {
		defer wg.Done()

		time.Sleep(3 * time.Second)
		gormDB.Debug().Create(&model.T{ID: 7, C: 7, D: 7}) // block
		common.PrintlnAllData(gormDB, "3")
	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}

func MysqlLock3() {
	gormDB := common.InitMysql()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		tx := gormDB.Begin()
		defer tx.Commit()

		tx.Debug().Raw("select * from ts where id >= 10 and id < 25 for update").Scan(&model.T{}) // [10,15]
		// tx.Debug().Raw("select * from ts where id = 10 for update").Scan(&model.T{}) // 10
		common.PrintlnAllData(tx, "1")

		time.Sleep(6 * time.Second)
	}()

	go func() {
		defer wg.Done()

		time.Sleep(2 * time.Second)
		gormDB.Debug().Create(&model.T{ID: 8, C: 8, D: 8})
		gormDB.Debug().Create(&model.T{ID: 13, C: 13, D: 13})

		common.PrintlnAllData(gormDB, "2")

	}()

	go func() {
		defer wg.Done()

		time.Sleep(3 * time.Second)

		// update d = d + 1 where id = 10
		//gormDB.Debug().Model(&model.T{}).Where("id = ?", 15).UpdateColumn("d", gorm.Expr("d + ?", 1))
		gormDB.Debug().Raw("update ts set d = d + 1 where id = 15").Scan(&model.T{})

		common.PrintlnAllData(gormDB, "3")
	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}
