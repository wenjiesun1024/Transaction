package main

import (
	mysql "Transaction/Mysql"
	pg "Transaction/Pg"

	_ "github.com/lib/pq"
)

func main() {
	MysqlTest := true

	if MysqlTest {
		// mysql.MysqlCurrentReadAndSnapRead()
		// mysql.MysqlDeadLock()
		mysql.MysqlPhantomRead()

		// mysql.MysqlLock4()
	} else {
		pg.PGCurrentReadAndSnapRead()

		// pg.PGPhantomRead()
		// pg.PGDeadLock()
		// pg.PGUpdate()
	}
}
