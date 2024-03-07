🍭 mysqldump
---
✅ 执行 mysqldump.sh 完成全量备份
变量 (必须):
1. BACKUP_DIR：备份目录
2. MYSQL_UNAME：用户
3. MYSQL_PWORD：密码
4. PATH：mysqldump命令路径
5. KEEP_BACKUPS_FOR：保留天数(默认7天)

🌐 恢复(示例)
```sh
for file in ./2024-03-07/*.sql.gz; do gunzip > "$file" | mysql -u your_username -p your_database_name; done
```
⛳️ cron(示例)

例：每天定时凌晨1点10分备份数据库

```
10 1 * * * /bin/bash /data/mysql/.sh/mysql_backup.sh
```

