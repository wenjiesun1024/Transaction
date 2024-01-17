package mysql

import (
	"Transaction/common"
	"Transaction/model"
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
		defer wg.Done()

		tx := gormDB.Begin()
		defer tx.Commit()

		common.PrintlnAllData(tx, "1")

		time.Sleep(3 * time.Second)

		// common.PrintlnAllData(tx, "3", clause.Locking{Strength: "UPDATE"})
		// common.PrintlnAllData(tx, "3", clause.Locking{Strength: "SHARE"})
		common.PrintlnAllData(tx, "3")

		tx.Model(&model.T{}).Where("id = ?", 5).UpdateColumn("d", gorm.Expr("d + ?", 10))

		common.PrintlnAllData(tx, "4")
	}()

	go func() {
		defer wg.Done()

		time.Sleep(2 * time.Second)

		tx := gormDB.Begin()
		defer tx.Commit()

		tx.Model(&model.T{}).Where("id = ?", 5).Update("d", 100000)
		tx.Model(&model.T{}).Create(&model.T{ID: 6, C: 6, D: 6})
		common.PrintlnAllData(tx, "2")
	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}
