package util

import (
	"math"
	"math/rand"
	"time"

	"DbSecurity/model"
	"DbSecurity/service"
	"github.com/shopspring/decimal"
)

type DifferentialPrivacy struct {
}

func (dp *DifferentialPrivacy) DifferPrivacy(patientsCopy *[]model.PatientCopy) {
	maxWeight, minWeight := dp.getMaxAndMin(patientsCopy, "weight")
	maxBill, minBill := dp.getMaxAndMin(patientsCopy, "bill")
	sensitivityWeight := maxWeight.Sub(minWeight)
	sensitivityBill := maxBill.Sub(minBill)
	epsilonWeight := decimal.NewFromInt(100)
	epsilonBill := decimal.NewFromInt(100)
	dp.laplace(sensitivityWeight, epsilonWeight, patientsCopy, "weight")
	dp.laplace(sensitivityBill, epsilonBill, patientsCopy, "bill")
}

//生成噪音
func (dp *DifferentialPrivacy) noiseCount(sensitivity decimal.Decimal, epsilon decimal.Decimal) decimal.Decimal {
	rand.Seed(time.Now().Unix())
	beta := sensitivity.Div(epsilon).Round(3)
	u1 := rand.Float64()
	u2 := rand.Float64()
	var nValue decimal.Decimal
	if u1 < 0.5 {
		nValue = beta.Neg().Mul(decimal.NewFromFloat(math.Log(1 - u2))).Round(1)
	} else {
		nValue = beta.Mul(decimal.NewFromFloat(u2)).Round(1)
	}
	return nValue
}

//laplace变换
func (dp *DifferentialPrivacy) laplace(sensitivity decimal.Decimal, epsilon decimal.Decimal, patientsCopy *[]model.PatientCopy,
	name string) {
	length := len(*patientsCopy)
	var patientCopyService service.PatientService
	if name == "weight" {
		for i := 0; i < length; i++ {
			(*patientsCopy)[i].Weight = (*patientsCopy)[i].Weight.Add(dp.noiseCount(sensitivity, epsilon))
			_ = patientCopyService.UpdatePatientCopy("weight", (*patientsCopy)[i].Weight, (*patientsCopy)[i].Id)
		}
	} else {
		for i := 0; i < length; i++ {
			(*patientsCopy)[i].Bill = (*patientsCopy)[i].Bill.Add(dp.noiseCount(sensitivity, epsilon))
			_ = patientCopyService.UpdatePatientCopy("bill", (*patientsCopy)[i].Bill, (*patientsCopy)[i].Id)
		}
	}
}

//获取最大值和最小值
func (dp *DifferentialPrivacy)getMaxAndMin(patientsCopy *[]model.PatientCopy, name string) (decimal.Decimal, decimal.Decimal) {
	length := len(*patientsCopy)
	var max, min decimal.Decimal
	if name == "weight" {
		max = (*patientsCopy)[0].Weight
		min = (*patientsCopy)[0].Weight
		for i := 0; i < length; i++ {
			if (*patientsCopy)[i].Weight.LessThan(min) {
				min = (*patientsCopy)[i].Weight
			}
			if max.LessThan((*patientsCopy)[i].Weight) {
				max = (*patientsCopy)[i].Weight
			}
		}
	} else {
		max = (*patientsCopy)[0].Bill
		min = (*patientsCopy)[0].Bill
		for i := 0; i < length; i++ {
			if (*patientsCopy)[i].Bill.LessThan(min) {
				min = (*patientsCopy)[i].Bill
			}
			if max.LessThan((*patientsCopy)[i].Bill) {
				max = (*patientsCopy)[i].Bill
			}
		}
	}
	return max, min
}