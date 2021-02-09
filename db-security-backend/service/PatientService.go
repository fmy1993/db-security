package service

import (
	"db-security-backend/dao"
	"db-security-backend/model"
)

type PatientService struct {
}

//得到patient表的所有数据
func (ps *PatientService) GetPatients() (*[]model.Patient, error) {
	var patientDao = dao.NewPatientDao()
	return patientDao.QueryAllPatients()
}

//判断克隆表是否存在
func (ps *PatientService) IsPatientCopyExist() (bool, error) {
	var patientDao = dao.NewPatientDao()
	return patientDao.IsPatientCopyExist()
}

//判断用户表是否存在
func (ps *PatientService) IsUserPatientCopyExist(phone string) (bool, error) {
	var patientDao = dao.NewPatientDao()
	return patientDao.IsUserPatientCopyExist(phone)
}

//创建克隆表
func (ps *PatientService) CreateNewPatientCopyTable() (bool, error) {
	var patientDao = dao.NewPatientDao()
	exist, err := patientDao.IsPatientCopyExist()
	if err != nil {
		return false, err
	}
	if exist {
		err = patientDao.DropPatientCopyTable()
		if err != nil {
			return false, err
		}
	}
	err = patientDao.CreatePatientCopyTable()
	if err != nil {
		return false, err
	}
	return true, nil
}

//创建用户表
func (ps *PatientService) CreateNewUserPatientCopyTable(phone string) (bool, error) {
	var patientDao = dao.NewPatientDao()
	exist, err := patientDao.IsUserPatientCopyExist(phone)
	if err != nil {
		return false, err
	}
	if exist {
		err = patientDao.DropUserPatientCopyTable(phone)
		if err != nil {
			return false, err
		}
	}
	err = patientDao.CreateUserPatientCopyTable(phone)
	if err != nil {
		return false, err
	}
	return true, nil
}

//获取克隆表所有数据
func (ps *PatientService) GetPatientCopyData() (*[]model.PatientCopy, error) {
	var patientDao = dao.NewPatientDao()
	return patientDao.QueryAllPatientCopy()
}

//获取用户克隆表所有数据
func (ps *PatientService) GetUserPatientCopyData(phone string) (*[]model.PatientCopy, error) {
	var patientDao = dao.NewPatientDao()
	return patientDao.QueryAllUserPatientCopy(phone)
}

//获取用户克隆表单页数据
func (ps *PatientService) GetPageUserPatientCopyData(phone string, pageSize int, page int) (*[]model.PatientCopy,
	error) {
	var patientDao = dao.NewPatientDao()
	return patientDao.QueryPageUserPatientCopy(phone, pageSize, page)
}

//获取原始表单页数据
func (ps *PatientService) GetPagePatientData(pageSize int, page int) (*[]model.Patient,
	error) {
	var patientDao = dao.NewPatientDao()
	return patientDao.QueryPagePatient(pageSize, page)
}

//条件获取用户克隆表单页数据
func (ps *PatientService) GetPageUserPatientCopyDataByPattern(phone string, pageSize int,
	page int, pattern string) (*[]model.PatientCopy,
	error) {
	var patientDao = dao.NewPatientDao()
	return patientDao.QueryPageUserPatientCopyLikePattern(phone, pageSize, page, pattern)
}

//条件获取原始表单页数据
func (ps *PatientService) GetPagePatientDataByPattern(pageSize int,
	page int, pattern string) (*[]model.Patient,
	error) {
	var patientDao = dao.NewPatientDao()
	return patientDao.QueryPagePatientLikePattern(pageSize, page, pattern)
}

//更新克隆表数据
func (ps *PatientService) UpdatePatientCopy(tupleName string, data interface{}, id int64) error {
	var patientDao = dao.NewPatientDao()
	return patientDao.UpdatePatientCopy(tupleName, data, id)
}

//更新用户表数据
func (ps *PatientService) UpdateUserPatientCopy(phone string, tupleName string, data interface{}, id int64) error {
	var patientDao = dao.NewPatientDao()
	return patientDao.UpdateUserPatientCopy(phone, tupleName, data, id)
}

//更新原始表数据
func (ps *PatientService) UpdatePatient(patient model.Patient) error {
	var patientDao = dao.NewPatientDao()
	return patientDao.UpdatePatient(patient)
}

//增加原始表数据
func (ps *PatientService) InsertPatient(patient model.Patient) error {
	var patientDao = dao.NewPatientDao()
	return patientDao.InsertPatient(patient)
}

//删除原始表数据
func (ps *PatientService) DeletePatient(patient model.Patient) error {
	var patientDao = dao.NewPatientDao()
	return patientDao.DeletePatient(patient)
}

//获取用户表某条数据
func (ps *PatientService) QueryOne(phone string, id int64) *model.PatientCopy{
	var patientDao = dao.NewPatientDao()
	return patientDao.QueryOne(phone, id)
}