package pg

import (
	"SQLIsolationLevelTest/common"
	"SQLIsolationLevelTest/model"
	"sync"
	"time"
)

func PostgresSQL() {
	gormDB := common.InitPG()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()
		// tx.Exec(`set transaction isolation level repeatable read`)

		time.Sleep(2 * time.Second)

		tx.Model(&model.User{}).Where("name = ?", "John").Update("age", 100000)
		common.PrintlnAllUsers(tx, "2")

		// time.Sleep(1 * time.Second)
		// gormDB.Model(&User{}).Create(&User{Name: "Tom", Age: 44, UserID: 3})
		time.Sleep(5 * time.Second)

	}()

	go func() {
		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()

		// gormDB.Model(&User{}).Create(&User{Name: "Tom", Age: 55, UserID: 3})
		// time.Sleep(5 * time.Second)

		// set transaction isolation level
		// tx.Exec(`set transaction isolation level repeatable read`)

		common.PrintlnAllUsers(tx, "1")

		tx.Model(&model.User{}).Create(&model.User{Name: "Tom", Age: 55, UserID: 3})

		time.Sleep(4 * time.Second)

		common.PrintlnAllUsers(tx, "3") // John 100000 (snap read)

		time.Sleep(4 * time.Second)
		// plus john's age 10
		// tx.Model(&User{}).Where("name = ?", "John").UpdateColumn("age", gorm.Expr("age + ?", 10))

		// delete john
		// tx.Model(&User{}).Where("name = ?", "John").Delete(&User{})

		// skip lock

		// err := tx.Model(&User{}).Where("name = ?", "John").Update("age", 200).Error // RC Block
		// if err != nil {
		// 	fmt.Println(err, "??")
		// }

		common.PrintlnAllUsers(tx, "4") // Johh 100010 (current read)

	}()

	wg.Wait()

	common.PrintlnAllUsers(gormDB, "end")
}
