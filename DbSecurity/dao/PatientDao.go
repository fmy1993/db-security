package dao

import (
	"DbSecurity/model"
	"DbSecurity/tool"
)

type PatientDao struct {
	*tool.Orm
}

func NewPatientDao() *PatientDao {
	return &PatientDao{tool.DbEngine}
}

//查询病人表中所有数据
func (pd *PatientDao) QueryAllPatients() (*[]model.Patient, error) {
	var patients []model.Patient
	err := pd.Find(&patients)
	if err != nil {
		return nil, err
	}
	return &patients, nil
}

//判断克隆表是否存在
func (pd *PatientDao) IsPatientCopyExist() (bool, error) {
	exist, err := pd.IsTableExist("patient_copy")
	return exist, err
}

//删除克隆表
func (pd *PatientDao) DropPatientCopyTable() error {
	err := pd.DropTables("patient_copy")
	return err
}

//创建克隆表
func (pd *PatientDao) CreatePatientCopyTable() error {
	_, err := pd.Exec("create table patient_copy select id, weight, high, age, phone, address, bill from patient;")
	if err != nil {
		return err
	}
	_, err = pd.Exec("alter table patient_copy modify id int primary key auto_increment;")
	return err
}

//获取克隆表所有数据
func (pd *PatientDao) QueryAllPatientCopy() (*[]model.PatientCopy, error) {
	var patientCopy []model.PatientCopy
	err := pd.SQL("select * from patient_copy").Find(&patientCopy)
	if err != nil {
		return nil, err
	}
	return &patientCopy, nil
}

//更新克隆表数据
func (pd *PatientDao) UpdatePatientCopy(tupleName string, data interface{}, id int64) error {
	_, err := pd.Exec("update patient_copy set "+tupleName+" = ? where id = ?;", data, id)
	if err != nil {
		return err
	}
	return nil
}

//更新原始表数据
func (pd *PatientDao) UpdatePatient(patient model.Patient) error {
	_, err := pd.ID(patient.Id).Update(patient)
	if err != nil {
		return err
	}
	return nil
}

//增加原始表数据
func (pd *PatientDao) InsertPatient(patient model.Patient) error {
	_, err := pd.InsertOne(&patient)
	if err != nil {
		return err
	}
	return nil
}

//删除原始表数据
func (pd *PatientDao) DeletePatient(patient model.Patient) error {
	_, err := pd.ID(patient.Id).Delete(patient)
	if err != nil {
		return err
	}
	return nil
}

//判断用户表是否存在
func (pd *PatientDao) IsUserPatientCopyExist(phone string) (bool, error) {
	exist, err := pd.IsTableExist("patient_copy" + phone)
	return exist, err
}

//删除用户表
func (pd *PatientDao) DropUserPatientCopyTable(phone string) error {
	err := pd.DropTables("patient_copy" + phone)
	return err
}

//创建用户表
func (pd *PatientDao) CreateUserPatientCopyTable(phone string) error {
	_, err := pd.Exec("create table patient_copy" + phone + " select * from patient_copy;")
	if err != nil {
		return err
	}
	_, err = pd.Exec("alter table patient_copy" + phone + " modify id int primary key auto_increment;")
	return err
}

//获取用户表所有数据
func (pd *PatientDao) QueryAllUserPatientCopy(phone string) (*[]model.PatientCopy, error) {
	var patientCopy []model.PatientCopy
	err := pd.SQL("select * from patient_copy" + phone).Find(&patientCopy)
	if err != nil {
		return nil, err
	}
	return &patientCopy, nil
}

//更新用户克隆表数据
func (pd *PatientDao) UpdateUserPatientCopy(phone string, tupleName string, data interface{}, id int64) error {
	_, err := pd.Exec("update patient_copy"+phone+" set "+tupleName+" = ? where id = ?;", data, id)
	if err != nil {
		return err
	}
	return nil
}

//获取用户表单页数据
func (pd *PatientDao) QueryPageUserPatientCopy(phone string, pageSize int, page int) (*[]model.PatientCopy, error) {
	var patientCopy []model.PatientCopy
	err := pd.SQL("select * from patient_copy"+phone+" limit ? offset ?", pageSize, (page-1)*pageSize).Find(&patientCopy)
	if err != nil {
		return nil, err
	}
	return &patientCopy, nil
}

//获取原始表单页数据
func (pd *PatientDao) QueryPagePatient(pageSize int, page int) (*[]model.Patient, error) {
	var patient []model.Patient
	err := pd.SQL("select * from patient limit ? offset ?", pageSize, (page-1)*pageSize).Find(&patient)
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

//条件获取用户表单页数据
func (pd *PatientDao) QueryPageUserPatientCopyLikePattern(phone string, pageSize int,
	page int, pattern string) (*[]model.PatientCopy, error) {
	var patientCopy []model.PatientCopy
	err := pd.SQL("select * from patient_copy"+phone+" where address like '%"+pattern+"%' limit ? offset ?", pageSize,
		(page-1)*pageSize).Find(&patientCopy)
	if err != nil {
		return nil, err
	}
	return &patientCopy, nil
}

//条件获取原始表单页数据
func (pd *PatientDao) QueryPagePatientLikePattern(pageSize int,
	page int, pattern string) (*[]model.Patient, error) {
	var patient []model.Patient
	err := pd.SQL("select * from patient where address like '%"+pattern+"%' limit ? offset ?", pageSize,
		(page-1)*pageSize).Find(&patient)
	if err != nil {
		return nil, err
	}
	return &patient, nil
}

//获取某条数据
func (pd *PatientDao) QueryOne(phone string, id int64) *model.PatientCopy{
	var patientCopy model.PatientCopy
	_, _ = pd.SQL("select * from patient_copy"+phone+" where id = ?", id).Get(&patientCopy)
	return &patientCopy
}
