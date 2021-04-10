package model

type Config struct {
	ConfigKey   string `xorm:"pk varchar(30)" json:"config_key"`
	ConfigValue string `xorm:"varchar(150)" json:"config_value"`
}
