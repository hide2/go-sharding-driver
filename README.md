## Sharding
- 简单轻量高效
- 使用命令行工具进行一键初始化分库分表
- 简单的扩容策略
- 分库分表和读写分离连接池
- 跨平台客户端Driver直连数据库，利于本机开发调试
- 无需额外部署和依赖，不需要安装和部署中间件，性能0损耗
- 简单易用的DB接口自动根据TableName/ShardingKey进行无感执行
- *自定义的Twitter snowflake style ID generator
- *Golang/PHP/Python等语言的Driver支持，兼容一定ORM组件

## Config
```
{
    "sharding_datasources": {
		"ds_0": {
    	    "write": "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
	        "read": "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
		},
		"ds_1": {
    	    "write": "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
	        "read": "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
		},
		"ds_2": {
    	    "write": "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
	        "read": "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
		},
		"ds_3": {
    	    "write": "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local",
	        "read": "root:root@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
		},
    },
    "sharding_tables": {
		"users": {
			"ddl": "sql/sharding/user.sql",
			"sharding_key": "fpid",
			"sharding_ds_num": 4,
			"sharding_table_num": 256
		},
		"user_bind_keys": {
			"ddl": "sql/sharding/user_bind_keys.sql",
			"sharding_key": "bind_key",
			"sharding_ds_num": 8,
			"sharding_table_num": 256
		},
	}
}
```

## QuickStart
```
# 一键分库分表
./sharding init_db

# 简单的扩容策略
分表=(sharding_key%sharding_table_num)
分库=(sharding_key/sharding_table_num%sharding_ds_num%len(sharding_datasources))
例如初始只启动有4个ds，但可以将sharding_ds_num设置为8，当发现4个ds容量不够时，只需将ds0～ds4对应建立从库ds5~ds7同步数据后修改配置文件加入新的sharding_datasources，重启服务即可

# DB接口-单Shard
user := User{}
db.Query(1001, "SELECT * FROM users WHERE fpid = ?", 1001).Scan(&user)
db.Exec(1001, "UPDATE users SET lang = ? WHERE fpid = ?", "en", 1001)
db.Exec("abcde", "INSERT INTO user_bind_keys(bind_key, fpid) VALUES(?, ?)", "abcde", 1001)

# DB接口-MultiShards
users := Users{}
db.MultiQuery("SELECT * FROM users WHERE level > 100 limit 10").Scan(&users)
```