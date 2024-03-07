package main

import (
	"os"

	"github.com/jarvanstack/mysqldump"
)

func main() {

	// 连接数据库
	dsn := "root:rootpasswd@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"

	f, _ := os.Create("dump.sql")

	_ = mysqldump.Dump(
		dsn,                     // 连接数据库信息
		mysqldump.WithData(),    // Option: // 导出表数据
		mysqldump.WithWriter(f), // Option: Writer (Default: os.Stdout)
	)
}
