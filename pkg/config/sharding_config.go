package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
)

// 配置文件
type ShardingConfigStruct struct {
	ShardingDatasources []ShardingDataSource `json:"sharding_datasources"`
	ShardingTables      []ShardingTable      `json:"sharding_tables"`
}
type ShardingDataSource struct {
	Name  string `json:"name"`
	Write string `json:"write"`
	Read  string `json:"read"`
}
type ShardingTable struct {
	Name             string `json:"name"`
	Ddl              string `json:"ddl"`
	ShardingKey      string `json:"sharding_key"`
	ShardingDsNum    int    `json:"sharding_ds_num"`
	ShardingTableNum int    `json:"sharding_table_num"`
}

var ShardingConfigFile string
var ShardingConfig *ShardingConfigStruct

func init() {
	// 解析命令行参数
	var shardingConfigFile string
	flag.StringVar(&shardingConfigFile, "sc", "config/sharding_config.json", "set configuration `file`")
	flag.Parse()

	// 加载配置文件sharding_config.json
	data, err := ioutil.ReadFile(shardingConfigFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	shardingConfig := ShardingConfigStruct{}
	err = json.Unmarshal(data, &shardingConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	ShardingConfigFile = shardingConfigFile
	ShardingConfig = &shardingConfig
}
