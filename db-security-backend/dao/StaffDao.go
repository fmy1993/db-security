package dao

import (
	"db-security-backend/config"
	"db-security-backend/model"
	"db-security-backend/param/staffParam"
)

type StaffDao struct {
	*config.Orm
}

func NewStaffDao() *StaffDao {
	return &StaffDao{config.DbEngine}
}

// GetAllStaff 获取所有员工数据
func (sd *StaffDao) GetAllStaff() *[]model.Staff {
	staffs := make([]model.Staff, 0)
	err := sd.Find(&staffs)
	if err != nil {
		return nil
	}
	return &staffs
}

// GetStaffByStaffId 通过staffId获取staff
func (sd *StaffDao) GetStaffByStaffId(staffId int64) *model.Staff {
	var staff model.Staff
	_, _ = sd.SQL("select * from staff where staff_id = ?", staffId).Get(&staff)
	return &staff
}

// GetPageStaffByPhone 获取单页原始表数据
func (sd *StaffDao) GetPageStaffByPhone(selectStaffParam staffParam.SelectStaffParam) (*[]model.Staff, interface{}, error) {
	var staffs []model.Staff
	err := sd.SQL("select * from staff where qualification like "+"'%"+selectStaffParam.Qualification+"%'"+" limit ? offset ?", 50, (selectStaffParam.Page-1)*50).Find(&staffs)
	sql := "select count(*) from staff where qualification like " + "'%" + selectStaffParam.Qualification + "%'"
	count, _ := sd.SQL(sql).Count()
	if err != nil {
		return nil, nil, err
	}
	return &staffs, count, nil
}

// InsertStaff 增加原始表数据
func (sd *StaffDao) InsertStaff(staff model.Staff) error {
	_, err := sd.InsertOne(&staff)
	if err != nil {
		return err
	}
	return nil
}

// DeleteStaff 删除原始表数据
func (sd *StaffDao) DeleteStaff(staff *model.Staff) error {
	_, err := sd.ID(staff.StaffId).Delete(staff)
	if err != nil {
		return err
	}
	return nil
}

// UpdateStaff 更新原始表数据
func (sd *StaffDao) UpdateStaff(staff model.Staff) error {
	_, err := sd.ID(staff.StaffId).Update(staff)
	if err != nil {
		return err
	}
	return nil
}
