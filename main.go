package main

import (
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	gorm.Model
	UserID int `gorm:"uniqueIndex;column:user_id"`
	Name   string
	Age    int
}

func printlnAllUsers(db *gorm.DB, tag string, Clauses ...clause.Expression) {
	var users []User
	db.Clauses(Clauses...).Find(&users)
	fmt.Println(tag)
	for _, user := range users {
		fmt.Println(user.Name, user.Age, user.UserID)
	}
	fmt.Println("------------------------")
}

func main() {
	//Mysql()
	// MySQLSlow()
	// PostgresSQL()
	MysqlDeadLock()
}

/*
current read: when we use insert、update、delete、select for update or select for share, we will use current read
*/
func Mysql() {
	gormDB := initMysql()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		time.Sleep(2 * time.Second)

		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()

		tx.Model(&User{}).Where("name = ?", "John").Update("age", 100000)
		tx.Model(&User{}).Create(&User{Name: "Tom", Age: 44, UserID: 3})
		printlnAllUsers(tx, "2")
	}()

	go func() {
		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()

		printlnAllUsers(tx, "1")

		time.Sleep(3 * time.Second)

		printlnAllUsers(tx, "3") // Johh 40 (snap read)
		// printlnAllUsers(tx, "3", clause.Locking{Strength: "UPDATE"}) // Johh 100000 (current read)
		// printlnAllUsers(tx, "3", clause.Locking{Strength: "SHARE"}) // Johh 100000 (current read)

		// plus john's age 10
		// must wait transaction 1 commit
		tx.Model(&User{}).Where("name = ?", "John").UpdateColumn("age", gorm.Expr("age + ?", 10))

		printlnAllUsers(tx, "4") // Johh 100010 (current read)
		/*
			为什么 John现在可见最新数据，而 Tom还是旧数据呢？
				https://juejin.cn/post/7134186501306318856

				update操作产生了当前读，那当前读肯定可以读到其他事务已经提了的数据，
				然后更新后产生一个新的 ReadView，这个新的 ReadView 为刚刚更新的数据，所以 John 可以读到最新的数据， 哪怕是 current read。

				但是 Tom 读到的还是旧数据，因为他的 ReadView 还是旧的，所以他读到的还是旧数据。
		*/

	}()

	wg.Wait()

	printlnAllUsers(gormDB, "end")
}

func PostgresSQL() {
	gormDB := initPG()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()
		// tx.Exec(`set transaction isolation level repeatable read`)

		time.Sleep(2 * time.Second)

		tx.Model(&User{}).Where("name = ?", "John").Update("age", 100000)
		printlnAllUsers(tx, "2")

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

		printlnAllUsers(tx, "1")

		tx.Model(&User{}).Create(&User{Name: "Tom", Age: 55, UserID: 3})

		time.Sleep(4 * time.Second)

		printlnAllUsers(tx, "3") // John 100000 (snap read)

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

		printlnAllUsers(tx, "4") // Johh 100010 (current read)

	}()

	wg.Wait()

	printlnAllUsers(gormDB, "end")
}

func MySQLSlow() {
	gormDB := initMysql()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		time.Sleep(2 * time.Second)

		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()

		//tx.Model(&User{}).Where("name = ?", "John").Update("age", 100000)
		tx.Model(&User{}).Where("user_id = ?", 2).Update("age", 100000)
		printlnAllUsers(tx, "4")

		// time.Sleep(1 * time.Second)
		// gormDB.Model(&User{}).Create(&User{Name: "Tom", Age: 44, UserID: 3})
		time.Sleep(2 * time.Second)
	}()

	go func() {
		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()

		// gormDB.Model(&User{}).Create(&User{Name: "Tom", Age: 55, UserID: 3})
		// time.Sleep(5 * time.Second)

		// set transaction isolation level
		// tx.Exec(`set transaction isolation level repeatable read`)

		printlnAllUsers(tx, "1")

		tx.Model(&User{}).Create(&User{Name: "Tom", Age: 55, UserID: 3})

		time.Sleep(4 * time.Second)

		printlnAllUsers(tx, "2") // John 100000 (snap read)

		// plus john's age 10
		tx.Model(&User{}).Where("name = ?", "John").UpdateColumn("age", gorm.Expr("age + ?", 10))

		// delete john
		// tx.Model(&User{}).Where("name = ?", "John").Delete(&User{})

		// skip lock

		// err := tx.Model(&User{}).Where("name = ?", "John").Update("age", 200).Error // RC Block
		// if err != nil {
		// 	fmt.Println(err, "??")
		// }

		printlnAllUsers(tx, "3") // Johh 100010 (current read)

	}()

	wg.Wait()

	printlnAllUsers(gormDB, "end")
}

func MysqlDeadLock() {
	gormDB := initMysql()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		time.Sleep(2 * time.Second)

		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()

		// --------------scenero 1-----------
		// 在 RR，for update 会锁住所有读过的行，所以这里会锁住所有的行， 同时所有区间上的间隙也会被锁住
		// 另外一个transaction 已经锁住所有的行，所以 这个transaction 会被阻塞直到另外一个transaction commit
		// printlnAllUsers(tx, "2", clause.Locking{Strength: "UPDATE"})
		// ----------------------------------

		// ~~~~~~~~~~~scenero 2~~~~~~~~~~~
		// 在 RR 情况下， 因为 user_id 上有索引，所以这里只会锁住 user_id = 4 的行
		// 但是由于间隙锁的原因，会有(2, +∞)的间隙锁
		// 注意间隙锁之间是相互不冲突的， 与它冲突的是 “往这个间隙里插入一个新行” 的操作
		// 所以另外一个transaction 会被阻塞。同理，这个transaction 也会被阻塞直到另外一个transaction commit
		// 这就导致了死锁
		tx.Model(&User{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id", 4).First(&User{})
		tx.Create(&User{Name: "Tom", Age: 44, UserID: 3})
		// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	}()

	go func() {
		tx := gormDB.Begin()
		defer tx.Commit()
		defer wg.Done()

		// --------------scenero 1-----------
		// printlnAllUsers(tx, "1", clause.Locking{Strength: "UPDATE"})
		// time.Sleep(10 * time.Second)
		// ----------------------------------

		// ~~~~~~~~~~~~scenero 2~~~~~~~~~~
		tx.Model(&User{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id", 4).First(&User{})
		time.Sleep(3 * time.Second)
		tx.Create(&User{Name: "Tom", Age: 44, UserID: 3})
		// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

	}()

	wg.Wait()

	printlnAllUsers(gormDB, "end")
}

func initMysql() *gorm.DB {
	// use gorm to connect mysql
	dsn := "root:my-secret-pw@tcp(127.0.0.1:3306)/your_database_name?charset=utf8mb4&parseTime=True&loc=Local"
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// use gorm to create a table
	gormDB.AutoMigrate(&User{})

	// clear all data
	gormDB.Unscoped().Where("1 = 1").Delete(&User{})

	// insert data
	gormDB.Create(&User{Name: "Bob", Age: 30, UserID: 1})
	gormDB.Create(&User{Name: "John", Age: 40, UserID: 2})
	// gormDB.Create(&User{Name: "John2", Age: 50, UserID: 12})

	return gormDB
}

func initPG() *gorm.DB {
	// use gorm to creat a table
	gormDB, err := gorm.Open(postgres.Open("host=localhost user=user password=pass dbname=postgres port=5432 sslmode=disable"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// set transaction isolation level
	err = gormDB.Exec(`set transaction isolation level repeatable read`).Error
	if err != nil {
		panic(err)
	}

	// use gorm to create a table
	gormDB.AutoMigrate(&User{})

	// clear all data
	gormDB.Unscoped().Where("1 = 1").Delete(&User{})
	// insert data
	gormDB.Create(&User{Name: "Bob", Age: 30, UserID: 1})
	gormDB.Create(&User{Name: "John", Age: 40, UserID: 2})

	return gormDB
}
