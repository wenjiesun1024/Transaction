package main

import (
	mysql "Transaction/Mysql"
	pg "Transaction/Pg"

	_ "github.com/lib/pq"
)

func main() {
	{
		// mysql.MysqlPhantomRead()
		// mysql.MysqlCurrentReadAndSnapRead()
		// mysql.MysqlDeadLock()
		mysql.MysqlLock4()
	}
	{
		// pg.PGPhantomRead()
		// pg.PGCurrentReadAndSnapRead()
		// pg.PGDeadLock()
		pg.PGUpdate()
	}
}
