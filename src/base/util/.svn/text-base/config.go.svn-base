package util

import (
	"base/log"
	"encoding/json"
	"io/ioutil"
)

/*
	自定义的配置解析管理器
	global为母版配置
	node如果重复则会覆盖母版
*/
type Config struct {
	filename  string
	node      string
	configMap map[string]interface{}
}

func NewConfig(filename string, node string) (*Config, error) {
	cfg := &Config{
		filename:  filename,
		node:      node,
		configMap: make(map[string]interface{}),
	}

	err := cfg.init()
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (cfg *Config) init() error {
	data, err := ioutil.ReadFile(cfg.filename)
	if err != nil {
		return err
	}

	dataMap := make(map[string]map[string]interface{})
	err = json.Unmarshal(data, &dataMap)
	if err != nil {
		log.Println(err)
		return err
	}

	if _, ok := dataMap["global"]; ok {
		cfg.parseNode(dataMap["global"])
	}

	if _, ok := dataMap[cfg.node]; ok {
		cfg.parseNode(dataMap[cfg.node])
	}

	log.Println(cfg.configMap)
	return nil
}

func (cfg *Config) parseNode(nodeMap map[string]interface{}) {

	for k, v := range nodeMap {
		cfg.configMap[k] = v
	}
}

func (cfg *Config) Get(name string) interface{} {

	if _, ok := cfg.configMap[name]; !ok {
		return nil
	}

	return cfg.configMap[name]
}

func (cfg *Config) Set(name string, value interface{}) {
	if _, ok := cfg.configMap[name]; ok {
		return
	}

	cfg.configMap[name] = value
}

func (cfg *Config) GetInt(name string) int {
	value := cfg.Get(name)
	if value == nil {
		return 0
	}

	//json normal is float64
	//so ....
	f, flag := value.(float64)
	if flag {
		return int(f)
	}

	n, _ := value.(int)

	return n
}

func (cfg *Config) GetString(name string) string {

	value := cfg.Get(name)
	if value == nil {
		return ""
	}

	str, flag := value.(string)
	if !flag {
		return ""
	}

	return str
}
