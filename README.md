# 准备工作
- 使用创建Mysql 和 PostgreSQL数据库
  `docker compose up -d`

# Mysql
## CurrentReadAndSnapRead
    使用 For Update 和 Share 的读是当前读，而普通的读是快照读。
    更新操作也是当前读

## MysqlDeadLock
    间隙锁之间是相互不冲突的， 与它冲突的是 “往这个间隙里插入一个新行” 的操作
    一个transaction 会被 deadlock rollback，所以还是会有一个transaction commit

## MysqlPhantomRead
    幻读是指一个事务在读取某一范围的数据时，另一个事务又在该范围内插入了新的行，当第一个事务再次读取该范围的数据时，会发现多了一些原本不存在的行。
    一开始使用快照读，没有读取到数据
    然后另外一个事务插入数据
    然后第一个事务更新所有数据，即使用当前读，读取到所有数据并提交
    然后再次读取数据，发现多了一些原本不存在的行，这就是幻读。