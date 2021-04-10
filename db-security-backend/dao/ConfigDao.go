package dao

import (
	"db-security-backend/config"
	"db-security-backend/model"
)

type ConfigDao struct {
	*config.Orm
}

func NewConfigDao() *ConfigDao {
	return &ConfigDao{config.DbEngine}
}

// GetConfigByColumn 获取config表中指定数据
func (cd *ConfigDao) GetConfigByColumn(column string) (*model.Config, error) {
	var cfg model.Config
	_, err := cd.Where("config_key = ?", column).Get(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
