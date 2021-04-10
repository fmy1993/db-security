package dao

import (
	"db-security-backend/config"
	"db-security-backend/model"
	"db-security-backend/param/staffParam"
	"log"
)

type StaffCopyDao struct {
	*config.Orm
}

func NewStaffCopyDao() *StaffCopyDao {
	return &StaffCopyDao{config.DbEngine}
}

// IsStaffCopyExist 检查克隆表是否存在，如果不存在就创建一个新的克隆表
func (scd *StaffCopyDao) IsStaffCopyExist() bool {
	exist, _ := scd.IsTableExist("staff_copy")
	if exist {
		_, err := scd.Exec("drop table staff_copy;")
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	_, err := scd.Exec("create table staff_copy select staff_id, height, weight, qualification, salary from staff;")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = scd.Exec("alter table staff_copy modify staff_id int primary key auto_increment;")
	return true
}

// IsUserStaffCopyExist 检查用户克隆表是否存在，如果不存在就创建一个新的用户克隆表
func (scd *StaffCopyDao) IsUserStaffCopyExist(phone string) bool {
	exist, _ := scd.IsTableExist("staff_copy_" + phone)
	if exist {
		_, err := scd.Exec("drop table staff_copy_" + phone + ";")
		if err != nil {
			log.Fatal(err.Error())
		}
	}
	_, err := scd.Exec("create table staff_copy_" + phone + " select staff_id, height, weight, qualification, salary from staff_copy;")
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = scd.Exec("alter table staff_copy_" + phone + " modify staff_id int primary key auto_increment;")
	return true
}

// UpdateCloneTable 更新克隆表数据
func (scd *StaffCopyDao) UpdateCloneTable(staffCopy *model.StaffCopy) int64 {
	affected, err := scd.ID(staffCopy.StaffId).Update(staffCopy)
	if err != nil {
		log.Fatal(err.Error())
	}
	return affected
}

// GetAllStaffCopy 获取所有克隆表数据
func (scd *StaffCopyDao) GetAllStaffCopy() *[]model.StaffCopy {
	staffCopys := make([]model.StaffCopy, 0)
	err := scd.Find(&staffCopys)
	if err != nil {
		return nil
	}
	return &staffCopys
}

// GetAllStaffCopyByPhone 获取所有用户表数据
func (scd *StaffCopyDao) GetAllStaffCopyByPhone(phone string) *[]model.StaffCopy {
	staffCopys := make([]model.StaffCopy, 0)
	res := scd.SQL("select * from staff_copy_" + phone)
	err := res.Find(&staffCopys)
	if err != nil {
		return nil
	}
	return &staffCopys
}

// GetPageStaffCopyByPhone 获取单页用户表数据
func (scd *StaffCopyDao) GetPageStaffCopyByPhone(phone string, selectStaffParam staffParam.SelectStaffParam) (*[]model.StaffCopy, interface{}, error) {
	var staffCopys []model.StaffCopy
	err := scd.SQL("select * from staff_copy_"+phone+" where qualification like "+"'%"+selectStaffParam.Qualification+"%'"+" limit ? offset ?", 50, (selectStaffParam.Page-1)*50).Find(&staffCopys)
	sql := "select count(*) from staff_copy_" + phone + " where qualification like " + "'%" + selectStaffParam.Qualification + "%'"
	count, _ := scd.SQL(sql).Count()
	if err != nil {
		return nil, nil, err
	}
	return &staffCopys, count, nil
}
