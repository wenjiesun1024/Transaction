package main

import (
	mysql "Transaction/Mysql"
	pg "Transaction/Pg"

	_ "github.com/lib/pq"
)

func main() {
	MysqlTest := false
	EnablePGRR := true

	if MysqlTest {
		// mysql.MysqlCurrentReadAndSnapRead()
		// mysql.MysqlDeadLock()
		mysql.MysqlPhantomRead()

		// mysql.MysqlLock4()
	} else {
		pg.PGCurrentReadAndSnapRead(EnablePGRR)
		// pg.PGPhantomRead(EnablePGRR)
		// pg.PGDeadLock(EnablePGRR)
		// pg.PGUpdate(EnablePGRR)
	}
}
