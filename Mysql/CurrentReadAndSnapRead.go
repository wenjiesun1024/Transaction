package mysql

import (
	"Transaction/common"
	"Transaction/model"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/*
current read: when we use insert、update、delete、select for update or select for share, we will use current read
*/
func MysqlCurrentReadAndSnapRead() {
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

		common.WaitFor(cond, 4)

		switch "Normal" {
		case "Normal":
			common.PrintlnAllData(tx, "3") // id=5, c=5, d=5 (snap read)
		case "Update":
			common.PrintlnAllData(tx, "3", clause.Locking{Strength: "UPDATE"}) // id=5,c=5,d=10000 (current read)
		case "Share":
			common.PrintlnAllData(tx, "3", clause.Locking{Strength: "SHARE"}) // id=5,c=5,d=10000 (current read)
		}

		// must wait transaction 2 commit
		// why? 另外一个更新操作把读过的行锁住了，所以这里会被阻塞。
		tx.Model(&model.T{}).Where("id = ?", 5).UpdateColumn("d", gorm.Expr("d + ?", 10))

		common.PrintlnAllData(tx, "4") // id=5, c=5, d=100010 (snap read)
		/*
			为什么 id=5现在可见最新数据，而 id=6还是旧数据呢？
				https://juejin.cn/post/7134186501306318856

				update操作产生了当前读，那当前读肯定可以读到其他事务已经提了的数据，
				然后更新后产生一个新的 ReadView，这个新的 ReadView 为刚刚更新的数据，所以 id=5 可以读到最新的数据， 哪怕是 current read。

				但是 id=6 读到的还是旧数据，因为他的 ReadView 还是旧的，所以他读到的还是旧数据。
		*/

	}()

	go func() {
		defer wg.Done()
		tx := gormDB.Begin() // Transaction 2
		defer tx.Commit()
		common.WaitFor(cond, 2)

		tx.Model(&model.T{}).Where("id = ?", 5).Update("d", 100000)
		tx.Model(&model.T{}).Create(&model.T{ID: 6, C: 6, D: 6})
		common.PrintlnAllData(tx, "2")

		common.WaitFor(cond, 3)
	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}
