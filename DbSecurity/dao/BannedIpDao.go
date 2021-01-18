package dao

import (
	"log"

	"DbSecurity/model"
	"DbSecurity/tool"
)

type BannedIpDao struct {
	*tool.Orm
}

func NewBannedIpDao() *BannedIpDao {
	return &BannedIpDao{tool.DbEngine}
}

func (bid *BannedIpDao) IpIsExist(ip string) bool {
	var bannedIp model.BannedIp
	exist, err := bid.Where("ip = ?", ip).Get(&bannedIp)
	if err != nil {
		log.Fatal(err.Error())
	}
	return exist
}

//插入ip
func (bid *BannedIpDao) InsertIp(ip model.BannedIp) int64 {
	result, err := bid.InsertOne(&ip)
	if err != nil {
		log.Fatal(err.Error())
		return 0
	}
	return result
}

//ip是否存在
func (bid *BannedIpDao) IsIpExist(ip string) *model.BannedIp {
	var bannedIp model.BannedIp
	_, err := bid.Where("ip = ?", ip).Get(&bannedIp)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &bannedIp
}

//删除ip
func (bid *BannedIpDao) DeleteIp(ip *model.BannedIp) error {
	_, err := bid.ID(ip.Id).Delete(ip)
	if err != nil {
		return err
	}
	return nil
}

//获取所有数据
func (bid *BannedIpDao) QueryAllIp() (*[]model.BannedIp, error) {
	var ip []model.BannedIp
	err := bid.Find(&ip)
	if err != nil {
		return nil, err
	}
	return &ip, nil
}
