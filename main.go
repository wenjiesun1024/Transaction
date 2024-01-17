package main

import (
	mysql "Transaction/Mysql"
	pg "Transaction/Pg"

	_ "github.com/lib/pq"
)

func main() {
	if false {
		// mysql.MysqlPhantomRead()
		// mysql.MysqlCurrentReadAndSnapRead()
		mysql.MysqlDeadLock()
		// mysql.MysqlLock4()
	}
	if true {
		// pg.PGPhantomRead()
		pg.PGCurrentReadAndSnapRead()
		// pg.PGDeadLock()
		// pg.PGUpdate()
	}
}
