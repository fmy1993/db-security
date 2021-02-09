package util

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"

	"db-security-backend/service"
	"db-security-backend/tool"
)

type CsvOperation struct {
}

func (cp *CsvOperation) GenerateCsv(phone string) (int, error) {
	var patientService service.PatientService
	var cfg = tool.GetConfig()
	patientsCopy, err := patientService.GetUserPatientCopyData(phone)
	var data [][]string
	var p []string
	p = append(p, "id")
	p = append(p, "phone")
	p = append(p, "weight")
	p = append(p, "high")
	p = append(p, "age")
	p = append(p, "address")
	p = append(p, "bill")
	data = append(data, p)
	var length = len(*(patientsCopy))
	for i := 0; i < length; i++ {
		var p []string
		p = append(p, strconv.FormatInt((*(patientsCopy))[i].Id, 10))
		p = append(p, (*(patientsCopy))[i].Phone)
		p = append(p, (*(patientsCopy))[i].Weight.String())
		p = append(p, strconv.Itoa((*(patientsCopy))[i].High))
		p = append(p, strconv.Itoa((*(patientsCopy))[i].Age))
		p = append(p, (*(patientsCopy))[i].Address)
		p = append(p, (*(patientsCopy))[i].Bill.String())
		data = append(data, p)
	}
	if err != nil {
		return 0, err
	}
	ExportCsv(cfg.Path.Form+"patient"+phone+".csv", data)
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
