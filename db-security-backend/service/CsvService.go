package service

import (
	"db-security-backend/config"
	"encoding/csv"
	"log"
	"os"
	"strconv"
)

type CsvService struct {
	scs StaffCopyService
}

// 生成csv文件
func (cs *CsvService) GenerateCsv(phone string) (int, error) {
	var cfg = config.GetConfig()
	staffCopys, err := cs.scs.GetAllStaffCopyDataByPhone(phone)
	var data [][]string
	var p []string
	p = append(p, "staff_id")
	p = append(p, "height")
	p = append(p, "weight")
	p = append(p, "qualification")
	p = append(p, "salary")
	data = append(data, p)
	var length = len(*(staffCopys))
	for i := 0; i < length; i++ {
		var p []string
		p = append(p, strconv.FormatInt((*(staffCopys))[i].StaffId, 10))
		p = append(p, (*(staffCopys))[i].Height.String())
		p = append(p, (*(staffCopys))[i].Weight.String())
		p = append(p, (*(staffCopys))[i].Qualification)
		p = append(p, (*(staffCopys))[i].Salary.String())
		data = append(data, p)
	}
	if err != nil {
		return 0, err
	}
	ExportCsv(cfg.Path.Form+"staff"+phone+".csv", data)
	return length, nil
}

func ExportCsv(filePath string, data [][]string) {
	fp, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer fp.Close()
	_, _ = fp.WriteString("\xEF\xBB\xBF")
	w := csv.NewWriter(fp)
	_ = w.WriteAll(data)
	w.Flush()
}
