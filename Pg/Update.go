package pg

import (
	"Transaction/common"
	"Transaction/model"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

// 自动检测更新
func PGUpdate(EnableRR bool) {
	gormDB := common.InitPG(EnableRR)

	wg := sync.WaitGroup{}
	wg.Add(2)

	cond := &common.MyCond{
		Key:  int32(1),
		Cond: sync.NewCond(new(sync.Mutex)),
	}

	go func() {
		defer wg.Done()

		fmt.Println("Transaction 1 Err: ", gormDB.Transaction(func(tx *gorm.DB) error {

			common.PrintlnAllData(tx, "1")

			common.WaitFor(cond, 1)

			common.WaitFor(cond, 4)

			// update
			if err := tx.Model(&model.T{}).Where("id = ?", 1).Update("d", 550).Error; err != nil {
				return err
			}

			common.PrintlnAllData(tx, "3")

			// update, abort
			if err := tx.Model(&model.T{}).Where("id = ?", 10).Update("d", 550).Error; err != nil {
				return err
			}

			common.PrintlnAllData(tx, "4")
			return nil
		}))
	}()

	go func() {
		defer wg.Done()
		common.WaitFor(cond, 2)

		if err := gormDB.Model(&model.T{}).Where("id = ?", 10).Update("d", 500).Error; err != nil {
			fmt.Println("Transaction 2 Err: ", err)
		}
		common.PrintlnAllData(gormDB, "2")

		common.WaitFor(cond, 3)
	}()

	wg.Wait()
	common.PrintlnAllData(gormDB, "end")
}
