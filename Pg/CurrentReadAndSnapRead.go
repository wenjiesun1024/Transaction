package pg

import (
	"Transaction/common"
	"Transaction/model"
	"fmt"
	"sync"

	"gorm.io/gorm"
	// "gorm.io/gorm/clause"
)

// FIXME: 不能稳定得到预期结果
func PGCurrentReadAndSnapRead(EnableRR bool) {
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

			// if err := common.PrintlnAllData(tx, "3", clause.Locking{Strength: "UPDATE"}); err != nil {
			// 	return err
			// } // abort
			// if err := common.PrintlnAllData(tx, "3", clause.Locking{Strength: "SHARE"}); err != nil {
			// 	return err
			// } // abort
			common.PrintlnAllData(tx, "3") // id=5, c=5, d=5

			// 自动检测更新丢失
			// update id = 10, ok
			if err := tx.Model(&model.T{}).Where("id = ?", 10).UpdateColumn("d", gorm.Expr("d + ?", 10)).Error; err != nil {
				return err
			}
			common.PrintlnAllData(tx, "4")
			// update id = 6, ok, no data, no update
			if err := tx.Model(&model.T{}).Where("id = ?", 6).UpdateColumn("d", gorm.Expr("d + ?", 10)).Error; err != nil {
				return err
			}
			common.PrintlnAllData(tx, "5")
			// update id = 5, abort
			if err := tx.Model(&model.T{}).Where("id = ?", 5).UpdateColumn("d", gorm.Expr("d + ?", 10)).Error; err != nil {
				return err
			}

			common.PrintlnAllData(tx, "6")
			return nil
		}))

	}()

	go func() {
		defer wg.Done()
		fmt.Println("Transaction 2 Err: ", gormDB.Transaction(func(tx *gorm.DB) error {

			common.WaitFor(cond, 2)

			if err := tx.Model(&model.T{}).Where("id = ?", 5).Update("d", 100000).Error; err != nil {
				return err
			}
			if err := tx.Model(&model.T{}).Create(&model.T{ID: 6, C: 6, D: 6, E: 6}).Error; err != nil {
				return err
			}
			common.PrintlnAllData(tx, "2")

			common.WaitFor(cond, 3)
			return nil
		}))
	}()

	wg.Wait()

	common.PrintlnAllData(gormDB, "end")
}
