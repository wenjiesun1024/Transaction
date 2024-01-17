package main

import (
	mysql "Transaction/Mysql"
	pg "Transaction/Pg"

	_ "github.com/lib/pq"
)

func main() {
	if true {
		// mysql.MysqlPhantomRead()
		// mysql.MysqlCurrentReadAndSnapRead()
		// mysql.MysqlDeadLock()
		mysql.MysqlLock7()
	}
	if false {
		// pg.PGPhantomRead()
		// pg.PGCurrentReadAndSnapRead()
		// pg.PGDeadLock()
		pg.PGUpdate()
	}
}
