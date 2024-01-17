package mysql

import (
	"Transaction/common"
	"Transaction/model"
	"sync"
	"time"

	"gorm.io/gorm/clause"
)

// 主键等值查询
func MysqlLock() {
	gormDB := common.InitMysql()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		tx := gormDB.Begin()
		defer tx.Commit()

		tx.Model(&model.T{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", 7).First(&model.T{}) // for update select, 间隙锁(5, 10)
		common.PrintlnAllData(tx, "1")

		time.Sleep(5 * time.Second)
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

// 非唯一索引的等值查询
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

		time.Sleep(5 * time.Second)
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

// 主键范围查询
func MysqlLock3() {
	gormDB := common.InitMysql()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		tx := gormDB.Begin()
		defer tx.Commit()

		tx.Debug().Raw("select * from ts where id >= 10 and id < 11 for update").Scan(&model.T{}) // [10,15)
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
		gormDB.Debug().Raw("update ts set d = d + 1 where id = 15").Scan(&model.T{}) // not block, mysql version >= 8.0.18

		common.PrintlnAllData(gormDB, "3")
	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}

// 非唯一索引的范围查询
func MysqlLock4() {
	gormDB := common.InitMysql()

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		tx := gormDB.Begin()
		defer tx.Commit()

		tx.Debug().Raw("select * from ts where c >= 10 and c < 11 lock in share mode").Scan(&model.T{}) // (5,15)

		common.PrintlnAllData(tx, "1")

		time.Sleep(6 * time.Second)

		common.PrintlnAllData(tx, "1")
	}()

	go func() {
		defer wg.Done()

		time.Sleep(2 * time.Second)
		gormDB.Debug().Create(&model.T{ID: 8, C: 8, D: 8})

		common.PrintlnAllData(gormDB, "2")

	}()

	go func() {
		defer wg.Done()

		time.Sleep(3 * time.Second)

		// update d = d + 1 where id = 10
		gormDB.Debug().Raw("update ts set d = d + 1 where id = 15").Scan(&model.T{}) // not block, mysql version >= 8.0.18

		common.PrintlnAllData(gormDB, "3")
	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}

// 非唯一索引存在等值情况
func MysqlLock5() {
	gormDB := common.InitMysql()
	gormDB.Debug().Create(&model.T{ID: 30, C: 10, D: 30})

	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		tx := gormDB.Begin()
		defer tx.Commit()

		// delete from ts where c = 10
		tx.Debug().Raw("delete from ts where c = 10").Scan(&model.T{}) // ((c=5,id=5), (c=15,id=15))
		common.PrintlnAllData(tx, "1")

		time.Sleep(6 * time.Second)

		//common.PrintlnAllData(tx, "1")
	}()

	go func() {
		defer wg.Done()

		time.Sleep(2 * time.Second)
		gormDB.Debug().Create(&model.T{ID: 12, C: 12, D: 12})

		common.PrintlnAllData(gormDB, "2")

	}()

	go func() {
		defer wg.Done()

		time.Sleep(3 * time.Second)
		gormDB.Debug().Raw("update ts set d = d + 1 where c = 15").Scan(&model.T{}) // not block
	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}

// 非唯一索引存在等值情况 limit
func MysqlLock6() {
	gormDB := common.InitMysql()
	gormDB.Debug().Create(&model.T{ID: 30, C: 10, D: 30})

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		tx := gormDB.Begin()
		defer tx.Commit()

		// delete from ts where c = 10
		tx.Debug().Raw("delete from ts where c = 10 limit 2").Scan(&model.T{}) // ((c=5,id=5), (c=30,id=10)]
		common.PrintlnAllData(tx, "1")

		time.Sleep(6 * time.Second)
	}()

	go func() {
		defer wg.Done()

		time.Sleep(2 * time.Second)
		gormDB.Debug().Create(&model.T{ID: 12, C: 12, D: 12})

		common.PrintlnAllData(gormDB, "2")

	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}

// dead lock
func MysqlLock7() {
	gormDB := common.InitMysql()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		tx := gormDB.Begin()
		defer tx.Commit()

		tx.Debug().Raw("select * from ts where c = 10 lock in share mode").Scan(&model.T{}) // (5, 15)

		time.Sleep(3 * time.Second)

		tx.Debug().Create(&model.T{ID: 8, C: 8, D: 8})
	}()

	go func() {
		defer wg.Done()

		time.Sleep(2 * time.Second)
		gormDB.Debug().Raw("update ts set d = d + 1 where c = 10").Scan(&model.T{}) // block, 先（5，10）， 后 10

		common.PrintlnAllData(gormDB, "2")

	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}
