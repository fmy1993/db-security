package dao

import (
	"db-security-backend/config"
	"log"

	"db-security-backend/model"
)

type BannedIpDao struct {
	*config.Orm
}

func NewBannedIpDao() *BannedIpDao {
	return &BannedIpDao{config.DbEngine}
}

func (bid *BannedIpDao) IpIsExist(ip string) bool {
	var bannedIp model.BannedIp
	exist, err := bid.Where("ip = ?", ip).Get(&bannedIp)
	if err != nil {
		log.Fatal(err.Error())
	}
	return exist
}

// InsertIp 插入ip
func (bid *BannedIpDao) InsertIp(ip model.BannedIp) int64 {
	result, err := bid.InsertOne(&ip)
	if err != nil {
		log.Fatal(err.Error())
		return 0
	}
	return result
}

// IsIpExist ip是否存在
func (bid *BannedIpDao) IsIpExist(ip string) *model.BannedIp {
	var bannedIp model.BannedIp
	_, err := bid.Where("ip = ?", ip).Get(&bannedIp)
	if err != nil {
		log.Fatal(err.Error())
	}
	return &bannedIp
}

// DeleteIp 删除ip
func (bid *BannedIpDao) DeleteIp(ip *model.BannedIp) error {
	_, err := bid.ID(ip.Id).Delete(ip)
	if err != nil {
		return err
	}
	return nil
}

// QueryAllIp 获取所有数据
func (bid *BannedIpDao) QueryAllIp() (*[]model.BannedIp, error) {
	var ip []model.BannedIp
	err := bid.Find(&ip)
	if err != nil {
		return nil, err
	}
	return &ip, nil
}

// GetIpByIpId 根据ipId获取ip
func (bid *BannedIpDao) GetIpByIpId(ipId int64) *model.BannedIp {
	var ip = model.BannedIp{}
	_, err := bid.Where("id = ?", ipId).Get(&ip)
	if err != nil {
		return nil
	}
	return &ip
}
