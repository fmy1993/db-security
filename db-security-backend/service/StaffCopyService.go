package service

import (
	"db-security-backend/dao"
	"db-security-backend/model"
	"db-security-backend/param"
	"db-security-backend/param/staffParam"
	"github.com/shopspring/decimal"
)

type StaffCopyService struct {
}

// UpdateCloneTable 更新克隆表数据
func (scs *StaffCopyService) UpdateCloneTable(staffCopy *model.StaffCopy) int64 {
	var staffCopyDao = dao.NewStaffCopyDao()
	affected := staffCopyDao.UpdateCloneTable(staffCopy)
	return affected
}

// GetAllStaffCopyDataByPhone 获取用户表数据
func (scs *StaffCopyService) GetAllStaffCopyDataByPhone(phone string) (*[]model.StaffCopy, error) {
	var staffCopyDao = dao.NewStaffCopyDao()
	var staffCopys *[]model.StaffCopy
	staffCopys = staffCopyDao.GetAllStaffCopyByPhone(phone)
	return staffCopys, nil
}

// GetAllStaffCopyData 获取克隆表数据
func (scs *StaffCopyService) GetAllStaffCopyData() (*[]model.StaffCopy, error) {
	var staffCopyDao = dao.NewStaffCopyDao()
	var staffCopys *[]model.StaffCopy
	staffCopys = staffCopyDao.GetAllStaffCopy()
	return staffCopys, nil
}

// IsUserStaffCopyExist 检查用户克隆表是否存在,不存在就创建
func (scs *StaffCopyService) IsUserStaffCopyExist(phone string) bool {
	var staffCopyDao = dao.NewStaffCopyDao()
	return staffCopyDao.IsUserStaffCopyExist(phone)
}

// GetPageStaffCopyByPhone 获取用户表单页数据
func (scs *StaffCopyService) GetPageStaffCopyByPhone(phone string, param staffParam.SelectStaffParam) (*[]model.StaffCopy, interface{}, error) {
	var staffCopyDao = dao.NewStaffCopyDao()
	return staffCopyDao.GetPageStaffCopyByPhone(phone, param)
}

func (scs *StaffCopyService) AnalysisStaff(data *[]model.StaffCopy, analysisParam param.AnalysisParam) map[string]int {
	length := len(*data)
	var response = make(map[string]int)
	if analysisParam.Salary != 0 && analysisParam.Weight == 0 && analysisParam.Height == 0 {
		var (
			salaryAbove = 0
			salaryBelow = 0
		)
		for i := 0; i < length; i++ {
			if (*data)[i].Salary.LessThan(decimal.NewFromInt(int64(analysisParam.Salary))) {
				salaryBelow += 1
			} else {
				salaryAbove += 1
			}
		}
		response["data_length"] = length
		response["salary_above"] = salaryAbove
		response["salary_below"] = salaryBelow
	} else if analysisParam.Weight != 0 && analysisParam.Salary == 0 && analysisParam.Height == 0 {
		var (
			weightAbove = 0
			weightBelow = 0
		)
		for i := 0; i < length; i++ {
			if (*data)[i].Weight.LessThan(decimal.NewFromInt(int64(analysisParam.Weight))) {
				weightBelow += 1
			} else {
				weightAbove += 1
			}
		}
		response["data_length"] = length
		response["weight_above"] = weightAbove
		response["weight_below"] = weightBelow
	} else if analysisParam.Weight != 0 && analysisParam.Height != 0 && analysisParam.Salary == 0 {
		var (
			weightAboveHeightAbove = 0
			weightBelowHeightAbove = 0
			weightAboveHeightBelow = 0
			weightBelowHeightBelow = 0
		)
		for i := 0; i < length; i++ {
			if decimal.NewFromInt(int64(analysisParam.Weight)).LessThan((*data)[i].Weight) && decimal.NewFromInt(int64(analysisParam.Height)).LessThanOrEqual((*data)[i].Height) {
				weightAboveHeightAbove += 1
			} else if decimal.NewFromInt(int64(analysisParam.Weight)).LessThan((*data)[i].Weight) && (*data)[i].
				Height.LessThan(decimal.NewFromInt(int64(analysisParam.Height))) {
				weightAboveHeightBelow += 1
			} else if (*data)[i].Weight.LessThan(decimal.NewFromInt(int64(analysisParam.Weight))) && decimal.NewFromInt(int64(analysisParam.Height)).LessThanOrEqual((*data)[i].Height) {
				weightBelowHeightAbove += 1
			} else {
				weightBelowHeightBelow += 1
			}
		}
		response["data_length"] = length
		response["weight_above_height_above"] = weightAboveHeightAbove
		response["weight_below_height_above"] = weightBelowHeightAbove
		response["weight_above_height_below"] = weightAboveHeightBelow
		response["weight_below_height_below"] = weightBelowHeightBelow
	} else if analysisParam.Weight == 0 && analysisParam.Salary == 0 && analysisParam.Height != 0 {
		var (
			heightAbove = 0
			heightBelow = 0
		)
		for i := 0; i < length; i++ {
			if decimal.NewFromInt(int64(analysisParam.Height)).LessThan((*data)[i].Height) {
				heightAbove += 1
			} else {
				heightBelow += 1
			}
		}
		response["data_length"] = length
		response["height_above"] = heightAbove
		response["height_below"] = heightBelow
	}
	return response
}
