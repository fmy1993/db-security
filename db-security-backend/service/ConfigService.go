package service

import "db-security-backend/dao"

type ConfigService struct {
}

func (cs *ConfigService) GetConfigValueByKey(key string) string {
	var configDao = dao.NewConfigDao()
	column, err := configDao.GetConfigByColumn(key)
	if err != nil {
		return ""
	}
	return column.ConfigValue
}
