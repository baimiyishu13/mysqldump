package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

/*
	备份数据库
	注意：这个示例假设`mysql`和`mysqldump`命令在你的`PATH`中可用。如果不是这样，可能需要提供完整的路径
	这个Go CLI工具接受四个命令行参数：
		1.`backupDir`（备份目录）
		2.`mysqlUname`（MySQL用户名）
		3.`mysqlPword`（MySQL密码）
		4.`keepBackupsFor`（保留备份的天数
	备份文件的路径格式为`<backupDir>/<YYYYMMDD>/<database>.sql.gz`。
	备份文件的父目录是`<backupDir>/<YYYYMMDD>`，其中`<YYYYMMDD>`是格式化后的当前日期
*/

func main() {

	// 定义命令行参数
	backupDir := flag.String("backupDir", "/data/backup", "备份目录")
	mysqlUname := flag.String("mysqlUname", "", "MySQL 用户名")
	mysqlPword := flag.String("mysqlPword", "", "MySQL 密码")
	keepBackupsFor := flag.Int("keepBackupsFor", 7, "保留备份的天数")
	flag.Parse()

	// 执行备份
	err := backupDatabases(*backupDir, *mysqlUname, *mysqlPword, *keepBackupsFor)
	if err != nil {
		fmt.Println("🤕 备份失败:", err)
		os.Exit(1)
	} else {
		fmt.Println("🎉 全量备份成功!")
	}
}

func backupDatabases(backupDir, mysqlUname, mysqlPword string, keepBackupsFor int) error {
	rmdir := backupDir
	cmd := exec.Command("find", rmdir, "-type", "d", "-ctime", fmt.Sprintf("+%d", keepBackupsFor), "-exec", "rm", "-rf", "{}", "\\;")
	fmt.Println(cmd)
	err := cmd.Run()

	//获取数据库列表
	//  /usr/local/bin/mysql -u root -p123 -e "SHOW DATABASES" |awk -F " " '{if (NR!=1) print $1}')
	cmd = exec.Command("bash", "-c", fmt.Sprintf("/usr/local/bin/mysql -u %s -p%s -e \"SHOW DATABASES\" | awk '{if (NR!=1) print $1}'", mysqlUname, mysqlPword))
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("🤕 获取数据库列表失败: %w", err)
	}
	databases := strings.Fields(string(output))
	if len(databases) == 0 {
		return fmt.Errorf("获取数据库列列表为空")
	} else {
		fmt.Println("👍 获取数据库列表成功")
	}

	// 备份每个数据库
	var currentTime = time.Now().Format("20060102")
	var wg sync.WaitGroup
	for _, database := range databases {
		if database == "Database" {
			continue
		}
		wg.Add(1)
		go func(db string) {
			defer wg.Done()
			backupFile := fmt.Sprintf("%s/%s/%s.sql.gz", backupDir, currentTime, db)
			backupFileDir := fmt.Sprintf("%s/%s/", backupDir, currentTime)
			// 获取备份文件的父目录
			backupDirPath := filepath.Dir(backupFileDir)
			// 创建备份文件的父目录
			err := os.MkdirAll(backupDirPath, 0755)
			if err != nil {
				fmt.Printf("🤕 创建备份目录 %s 失败: %s\n", backupDirPath, err.Error())
				return
			}
			cmdStr := fmt.Sprintf("mysqldump -u %s -p%s %s | gzip -9 > %s", mysqlUname, mysqlPword, db, backupFile)
			cmd := exec.Command("bash", "-c", cmdStr)
			//fmt.Println(cmd)
			err = cmd.Run()
			if err != nil {
				fmt.Printf("🤕 备份数据库 %s 失败: %s\n", db, err.Error())
			} else {
				fmt.Printf("👍 备份数据库： %s 成功\n", db)
			}
		}(database)
	}
	wg.Wait()
	return nil
}
