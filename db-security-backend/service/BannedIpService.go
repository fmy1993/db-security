package service

import (
	"db-security-backend/dao"
	"db-security-backend/model"
)

type BannedIpService struct {
}

//ip是否存在
func (bis *BannedIpService) IpExist(ip string) bool {
	bannedIpDao := dao.NewBannedIpDao()
	return bannedIpDao.IpIsExist(ip)
}

//封禁ip
func (bis *BannedIpService) BanIp(ip string) {
	bannedIpDao := dao.NewBannedIpDao()
	if bannedIpDao.IsIpExist(ip).Id == 0 {
		var bannedIp model.BannedIp
		bannedIp.Ip = ip
		bannedIpDao.InsertIp(bannedIp)
	}
}

//解封ip
func (bis *BannedIpService) FreeIp(ip *model.BannedIp) error {
	bannedIpDao := dao.NewBannedIpDao()
	return bannedIpDao.DeleteIp(ip)
}

//获取所有数据
func (bis *BannedIpService) GetAllIp()(*[]model.BannedIp, error) {
	var ipDao = dao.NewBannedIpDao()
	return ipDao.QueryAllIp()
}