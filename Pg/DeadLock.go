package pg

import (
	"Transaction/common"
	"Transaction/model"
	"fmt"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func PGDeadLock(EnableRR bool) {
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

			tx.Model(&model.T{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id", 3).First(&model.T{})

			common.WaitFor(cond, 4)

			if err := tx.Create(&model.T{ID: 3, C: 3, D: 3, E: 3}).Error; err != nil {
				return err
			}

			return nil
		}))
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Transaction 2 Err: ", gormDB.Transaction(func(tx *gorm.DB) error {

			common.WaitFor(cond, 2)

			common.PrintlnAllData(tx, "2")

			tx.Model(&model.T{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id", 4).First(&model.T{})

			common.WaitFor(cond, 3)

			if err := tx.Create(&model.T{ID: 4, C: 4, D: 4, E: 4}).Error; err != nil {
				return err
			}

			return nil
		}))
	}()

	wg.Wait()

	// 其中一个transaction 会被 deadlock rollback，所以还是会有一个transaction commit
	common.PrintlnAllData(gormDB, "end")
}
