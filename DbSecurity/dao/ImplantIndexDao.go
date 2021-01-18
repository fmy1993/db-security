package dao

import (
	"DbSecurity/model"
	"DbSecurity/tool"
)

type ImplantIndexDao struct {
	*tool.Orm
}

//初始化
func NewImplantIndexDao() *ImplantIndexDao {
	return &ImplantIndexDao{tool.DbEngine}
}

//增加数据
func (iid *ImplantIndexDao) CreateRecord(index model.ImplantIndex) (int64, error) {
	res, err := iid.InsertOne(&index)
	if err != nil {
		return 0, err
	}
	return res, nil
}

//获取所有数据
func (iid *ImplantIndexDao) QueryAllII() (*[]model.ImplantIndex, error) {
	var ii []model.ImplantIndex
	err := iid.Find(&ii)
	if err != nil {
		return nil, err
	}
	return &ii, nil
}
