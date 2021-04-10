package service

import (
	"db-security-backend/dao"
	"db-security-backend/model"
	"db-security-backend/param/staffParam"
)

type StaffService struct {
}

// GetAllStaff 获取所有员工数据
func (ss *StaffService) GetAllStaff() (*[]model.Staff, error) {
	var staffDao = dao.NewStaffDao()
	var staffs *[]model.Staff
	staffs = staffDao.GetAllStaff()
	return staffs, nil
}

// GetPageStaff 获取原始表单页数据
func (ss *StaffService) GetPageStaff(param staffParam.SelectStaffParam) (*[]model.Staff, interface{}, error) {
	var staffDao = dao.NewStaffDao()
	return staffDao.GetPageStaffByPhone(param)
}

// InsertStaff 增加原始表数据
func (ss *StaffService) InsertStaff(staff model.Staff) error {
	var staffDao = dao.NewStaffDao()
	return staffDao.InsertStaff(staff)
}

// DeleteStaff 删除原始表数据
func (ss *StaffService) DeleteStaff(staffId int64) error {
	var staffDao = dao.NewStaffDao()
	var staff = staffDao.GetStaffByStaffId(staffId)
	return staffDao.DeleteStaff(staff)
}

// UpdateStaff 更新原始表数据
func (ss *StaffService) UpdateStaff(staff model.Staff) error {
	var staffDao = dao.NewStaffDao()
	return staffDao.UpdateStaff(staff)
}
