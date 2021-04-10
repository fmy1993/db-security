package service

import (
	"db-security-backend/config"
	"db-security-backend/dao"
	"db-security-backend/model"
	"github.com/shopspring/decimal"
	"log"
	"math"
	"math/rand"
	"time"
)

type DifferentialPrivacyService struct {
	ss  StaffService
	scs StaffCopyService
}

// DifferentialPrivacy 进行差分隐私
func (dps *DifferentialPrivacyService) DifferentialPrivacy() {
	maxSalary, minSalary := dps.getMaxAndMin()
	sensitivitySalary := maxSalary.Sub(minSalary)
	epsilonSalary := decimal.NewFromInt(100)
	var staffCopyDao = dao.NewStaffCopyDao()
	staffCopyDao.IsStaffCopyExist()
	dps.laplace(sensitivitySalary, epsilonSalary, staffCopyDao.GetAllStaffCopy())
}

// ExpMechanism 进行指数机制扰动, 返回最大身份证开头地区
func (dps *DifferentialPrivacyService) ExpMechanism(sensitivity decimal.Decimal) string {
	var staffDao = dao.NewStaffDao()
	var staffs = staffDao.GetAllStaff()
	idCardMap := make(map[string]int64, 0)
	for _, value := range *staffs {
		if _, ok := idCardMap[value.IdCard[:2]]; ok {
			idCardMap[value.IdCard[:2]] += 1
		} else {
			idCardMap[value.IdCard[:2]] = 1
		}
	}
	var sum = decimal.NewFromFloat(1)
	idCardProbability := make(map[string]decimal.Decimal, 0)
	var max int64 = 0
	var min int64 = math.MaxInt64
	for _, value := range idCardMap {
		if value > max {
			max = value
		}
		if value < min {
			min = value
		}
	}
	for key, value := range idCardMap {
		idCardProbability[key] = decimal.NewFromFloat(math.E).Pow(sensitivity.Mul(decimal.NewFromInt(value)).Div(decimal.NewFromInt(min)))
		sum = sum.Add(idCardProbability[key])
	}
	var idCardRange = make(map[string]struct {
		first  decimal.Decimal
		second decimal.Decimal
	}, 0)
	var start = decimal.NewFromFloat(0)
	for key, value := range idCardProbability {
		idCardProbability[key] = value.Div(sum)
		idCardRange[key] = struct {
			first  decimal.Decimal
			second decimal.Decimal
		}{first: start, second: start.Add(idCardProbability[key])}
		start = idCardRange[key].second
	}
	var ranNum = decimal.NewFromFloat(rand.Float64())
	for key, value := range idCardRange {
		if ranNum.LessThanOrEqual(value.second) && value.first.LessThanOrEqual(ranNum) {
			var configDao = dao.NewConfigDao()
			pro, _ := configDao.GetConfigByColumn(key)
			return pro.ConfigValue
		}
	}
	return "00"
}

// 进行laplace扰动
func (dps *DifferentialPrivacyService) laplace(sensitivity decimal.Decimal, epsilon decimal.Decimal, staffCopy *[]model.StaffCopy) {
	length := len(*staffCopy)
	session := config.DbEngine.NewSession()
	defer session.Close()
	_ = session.Begin()
	for i := 0; i < length; i++ {
		(*staffCopy)[i].Salary = (*staffCopy)[i].Salary.Add(dps.getNoise(sensitivity, epsilon))
		// dps.scs.UpdateCloneTable(&(*staffCopy)[i])
		_, _ = session.ID((*staffCopy)[i].StaffId).Update((*staffCopy)[i])
	}
	_ = session.Commit()
}

// 获取满足laplace分布的噪音
func (dps *DifferentialPrivacyService) getNoise(sensitivity decimal.Decimal, epsilon decimal.Decimal) decimal.Decimal {
	rand.Seed(time.Now().Unix())
	beta := sensitivity.Div(epsilon).Round(3)
	u1 := rand.Float64()
	u2 := rand.Float64()
	var noise decimal.Decimal
	if u1 < 0.5 {
		noise = beta.Neg().Mul(decimal.NewFromFloat(math.Log(1 - u2))).Round(1)
	} else {
		noise = beta.Mul(decimal.NewFromFloat(u2)).Round(1)
	}
	return noise
}

// 获取薪资的最大值和最小值
func (dps *DifferentialPrivacyService) getMaxAndMin() (decimal.Decimal, decimal.Decimal) {
	var max, min decimal.Decimal
	staffs, err := dps.ss.GetAllStaff()
	if err != nil {
		log.Fatal(err.Error())
	}
	max = (*staffs)[0].Salary
	min = (*staffs)[0].Salary
	for i := 0; i < len(*staffs); i++ {
		if (*staffs)[i].Salary.LessThan(min) {
			min = (*staffs)[i].Salary
		}
		if max.LessThan((*staffs)[i].Salary) {
			max = (*staffs)[i].Salary
		}
	}
	return max, min
}
