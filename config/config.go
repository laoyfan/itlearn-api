package config

import (
	"encoding/json"
	"github.com/go-ini/ini"
	"io/ioutil"
)

//结构体标签的多个键值对之间，必须用空格分割。，不能用逗号！！！，不能用逗号！！！，不能用逗号！！！

// 应用的配置结构体
type AppConfig struct {
	*ServerConfig `json:"server" ini:"server"`
	*MySQLConfig  `json:"mysql" ini:"mysql"`
	*RedisConfig  `json:"redis" ini:"redis"`
	*LogConfig    `json:"log" ini:"log"`
}

// web server配置
type ServerConfig struct {
	Port int `json:"port" ini:"port"`
}

// MySQL数据库配置
type MySQLConfig struct {
	Host     string `json:"host" ini:"host"`
	Username string `json:"username" ini:"username"`
	Password string `json:"password" ini:"password"`
	Port     int    `json:"port" ini:"port"`
	DB       string `json:"db" ini:"db"`
}

// redis配置
type RedisConfig struct {
	Host     string `json:"host" ini:"host"`
	Password string `json:"password" ini:"password"`
	Port     int    `json:"port" ini:"port"`
	DB       int    `json:"db" ini:"db"`
}

// Log配置
type LogConfig struct {
	Level      string `json:"level" ini:"level"`
	Filename   string `json:"filename" ini:"filename"`
	MaxSize    int    `json:"maxsize" ini:"maxsize"`
	MaxAge     int    `json:"max_age" ini:"max_age"`
	MaxBackups int    `json:"max_backups" ini:"max_backups"`
}

var Conf = new(AppConfig) // 定义了全局的配置文件实例

// Init 初始化
func Init(file string) error {
	jsonData, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(jsonData, Conf); err != nil {
		return err
	}
	return nil
}

func InitFromStr(str string) error {
	if err := json.Unmarshal([]byte(str), Conf); err != nil {
		return err
	}
	return nil
}

func InitFromIni(filename string) error {
	err := ini.MapTo(Conf, filename)
	if err != nil {
		panic(err)
	}
	return err
}
