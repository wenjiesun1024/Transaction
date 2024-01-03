package main

import (

	// pg "SQLIsolationLevelTest/Pg"

	mysql "SQLIsolationLevelTest/Mysql"

	_ "github.com/lib/pq"
)

func main() {
	{
		// mysql.MysqlCurrentReadAndSnapRead()
		mysql.MysqlDeadLock()
		// mysql.MysqlLock3()
	}
	{
		// pg.PostgresSQL()
	}
}

// func MysqlDeadLock2() {
// 	gormDB := initMysql()

// 	wg := sync.WaitGroup{}
// 	wg.Add(2)

// 	go func() {
// 		time.Sleep(2 * time.Second)

// 		tx := gormDB.Begin()
// 		defer tx.Commit()
// 		defer wg.Done()

// 		tx.Model(&User{}).Where("name = ?", "John").Update("age", 100000) //加了几个才被锁住？
// 		//tx.Model(&User{}).Where("user_id = ?", 2).Update("age", 100000)
// 		printlnAllUsers(tx, "4")

// 		// time.Sleep(1 * time.Second)
// 		// gormDB.Model(&User{}).Create(&User{Name: "Tom", Age: 44, UserID: 3})
// 		time.Sleep(2 * time.Second)
// 	}()

// 	go func() {
// 		tx := gormDB.Begin()
// 		defer tx.Commit()
// 		defer wg.Done()

// 		// gormDB.Model(&User{}).Create(&User{Name: "Tom", Age: 55, UserID: 3})
// 		// time.Sleep(5 * time.Second)

// 		// set transaction isolation level
// 		// tx.Exec(`set transaction isolation level repeatable read`)

// 		printlnAllUsers(tx, "1")

// 		tx.Model(&User{}).Create(&User{Name: "Tom", Age: 55, UserID: 3})

// 		time.Sleep(4 * time.Second)

// 		printlnAllUsers(tx, "2") // John 100000 (snap read)

// 		// plus john's age 10
// 		//tx.Model(&User{}).Where("name = ?", "John").UpdateColumn("age", gorm.Expr("age + ?", 10)) // 加锁全部
// 		tx.Model(&User{}).Where("user_id", 3).UpdateColumn("age", gorm.Expr("age + ?", 10)) // 加锁一个

// 		// delete john
// 		// tx.Model(&User{}).Where("name = ?", "John").Delete(&User{})

// 		// skip lock

// 		// err := tx.Model(&User{}).Where("name = ?", "John").Update("age", 200).Error // RC Block
// 		// if err != nil {
// 		// 	fmt.Println(err, "??")
// 		// }

// 		printlnAllUsers(tx, "3") // Johh 100010 (current read)

// 	}()

// 	wg.Wait()

// 	printlnAllUsers(gormDB, "end")
// }
