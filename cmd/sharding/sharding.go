package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	. "server/pkg/config"
	. "server/pkg/db"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("\nUsage: ./sharding init_db\n")
	} else {
		a := os.Args[1]
		if a == "init_db" {
			fmt.Println("\n[ShardingDBs]", ShardingDBs)
			// 执行ddl
			for _, t := range ShardingConfig.ShardingTables {
				data, err := ioutil.ReadFile(t.Ddl)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println("\n-------------------------- ddl:", t.Ddl)
				sqldata := string(data)
				sqls := strings.Split(sqldata, ";")
				for ds, sdb := range ShardingDBs {
					fmt.Println("\n----------->", ds)
					for i := 0; i < t.ShardingTableNum; i++ {
						suffix := "_" + strconv.FormatInt(int64(i), 10)
						for _, sql := range sqls {
							sql = strings.Trim(sql, " ")
							if sql != "" {
								// DROP TABLE
								if strings.Contains(sql, "DROP TABLE IF EXISTS") {
									sql = sql + suffix
								}
								// CREATE TABLE
								if strings.Contains(sql, "CREATE TABLE") {
									re := regexp.MustCompile(`CREATE TABLE (.*) `)
									sql = re.ReplaceAllString(sql, "CREATE TABLE ${1}"+suffix+" ")
								}
								fmt.Println(sql)
								sdb.Exec(sql)
							}
						}
					}
				}
			}
		} else {
			fmt.Println("\nUsage: ./sharding init_db\n")
		}
	}
}
