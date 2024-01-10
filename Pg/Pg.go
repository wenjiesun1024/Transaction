package pg

import (
	"Transaction/common"
	"Transaction/model"
	"sync"
	"time"
	// "gorm.io/gorm/clause"
)

func PostgresSQL() {
	gormDB := common.InitPG()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		tx := gormDB.Begin()
		defer wg.Done()
		defer tx.Commit()
		// tx.Exec(`set transaction isolation level repeatable read`)

		time.Sleep(2 * time.Second)

		tx.Model(&model.T{}).Where("name = ?", "John").Update("age", 100000)
		common.PrintlnAllData(tx, "2")

		// time.Sleep(1 * time.Second)
		// gormDB.Model(&User{}).Create(&User{Name: "Tom", Age: 44, UserID: 3})
		time.Sleep(5 * time.Second)

	}()

	go func() {
		tx := gormDB.Begin()
		defer wg.Done()
		defer tx.Commit()

		// gormDB.Model(&User{}).Create(&User{Name: "Tom", Age: 55, UserID: 3})
		// time.Sleep(5 * time.Second)

		// set transaction isolation level
		// tx.Exec(`set transaction isolation level repeatable read`)

		common.PrintlnAllData(tx, "1")

		// tx.Model(&model.User{}).Create(&model.User{Name: "Tom", Age: 55, UserID: 3})

		time.Sleep(4 * time.Second)

		common.PrintlnAllData(tx, "3") // John 100000 (snap read)

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

		common.PrintlnAllData(tx, "4") // Johh 100010 (current read)

	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}
