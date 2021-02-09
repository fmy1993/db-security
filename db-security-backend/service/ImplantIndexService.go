package service

import (
	"db-security-backend/dao"
	"db-security-backend/model"
)

type ImplantIndexService struct {
}

//添加数据
func (iis *ImplantIndexService) CreateRecord(index model.ImplantIndex) (int64, error) {
	var implantIndexDao = dao.NewImplantIndexDao()
	res, err := implantIndexDao.CreateRecord(index)
	if err != nil {
		return 0, err
	}
	return res, nil
}

//获取所有数据
func (iis *ImplantIndexService) GetAllII()(*[]model.ImplantIndex, error) {
	var iiDao = dao.NewImplantIndexDao()
	return iiDao.QueryAllII()
}