package mysql

import (
	"SQLIsolationLevelTest/common"
	"SQLIsolationLevelTest/model"
	"sync"
	"time"

	"gorm.io/gorm"
)

/*
current read: when we use insert、update、delete、select for update or select for share, we will use current read
*/
func MysqlCurrentReadAndSnapRead() {
	gormDB := common.InitMysql()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()

		common.PrintlnAllUsers(tx, "1")

		time.Sleep(3 * time.Second)

		common.PrintlnAllUsers(tx, "3") // Johh 40 (snap read)
		// printlnAllUsers(tx, "3", clause.Locking{Strength: "UPDATE"}) // Johh 100000 (current read)
		// printlnAllUsers(tx, "3", clause.Locking{Strength: "SHARE"}) // Johh 100000 (current read)

		// must wait transaction 1 commit
		// why? 另外一个更新操作把读过的行锁住了，所以这里会被阻塞。
		tx.Model(&model.User{}).Where("name = ?", "John").UpdateColumn("age", gorm.Expr("age + ?", 10))

		common.PrintlnAllUsers(tx, "4") // Johh 100010 (current read)
		/*
			为什么 John现在可见最新数据，而 Tom还是旧数据呢？
				https://juejin.cn/post/7134186501306318856

				update操作产生了当前读，那当前读肯定可以读到其他事务已经提了的数据，
				然后更新后产生一个新的 ReadView，这个新的 ReadView 为刚刚更新的数据，所以 John 可以读到最新的数据， 哪怕是 current read。

				但是 Tom 读到的还是旧数据，因为他的 ReadView 还是旧的，所以他读到的还是旧数据。
		*/

	}()

	go func() {
		time.Sleep(2 * time.Second)

		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()

		tx.Model(&model.User{}).Where("name = ?", "John").Update("age", 100000)
		tx.Model(&model.User{}).Create(&model.User{Name: "Tom", Age: 44, UserID: 3})
		common.PrintlnAllUsers(tx, "2")
	}()

	wg.Wait()

	common.PrintlnAllUsers(gormDB, "end")
}
