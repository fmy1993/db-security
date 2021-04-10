package config

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	AppName  string         `json:"app_name"`
	AppMode  string         `json:"app_mode"`
	AppHost  string         `json:"app_host"`
	AppPort  string         `json:"app_port"`
	Database DatabaseConfig `json:"database"`
	Redis    RedisConfig    `json:"redis"`
	Register RegisterConfig `json:"register"`
	Path     PathConfig     `json:"path"`
}

type DatabaseConfig struct {
	Driver   string `json:"driver"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	DbName   string `json:"db_name"`
	CharSet  string `json:"charset"`
	ShowSql  bool   `json:"show_sql"`
}

type RedisConfig struct {
	Addr        string `json:"addr"`
	Port        string `json:"port"`
	Password    string `json:"password"`
	Db          int    `json:"db"`
	MaxIdle     int    `json:"max_idle"`
	MaxActive   int    `json:"max_active"`
	IdleTimeout int    `json:"idle_timeout"`
}

type RegisterConfig struct {
	CheckCode string `json:"check_code"`
}

type PathConfig struct {
	Watermark string `json:"watermark"`
	ResForm   string `json:"res_form"`
	Res       string `json:"res"`
	PickUp    string `json:"pick_up"`
	OldIndex  string `json:"old_index"`
	Index     string `json:"index"`
	Form      string `json:"form"`
}

var _cfg *Config = nil

func GetConfig() *Config {
	return _cfg
}

func ParseConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&_cfg); err != nil {
		log.Fatal(err.Error())
	}
	return _cfg, err
}
