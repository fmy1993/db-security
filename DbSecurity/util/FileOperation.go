package util

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"DbSecurity/model"
	"DbSecurity/service"
	"DbSecurity/tool"
	"github.com/shopspring/decimal"
	"gocv.io/x/gocv"
)

type IndexResult struct {
	Id    string `json:"id"`
	Tuple Tuple  `json:"tuple"`
}

type Tuple struct {
	TupleId   string `json:"tuple_id"`
	TupleName string `json:"tuple_name"`
}

type IndexResultSlice struct {
	Result []IndexResult `json:"result"`
}

type DicTuple struct {
	Weight decimal.Decimal
	Bill   decimal.Decimal
}

//判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//得到索引表
func GenerateIndex() error {
	watermark, err := getWatermarkSlice("old_pic.bmp")
	if err != nil {
		return err
	}
	patientService := service.PatientService{}
	patients, err := patientService.GetPatients()
	if err != nil {
		return err
	}
	embedIdQueue := getIndexStruct(watermark, *patients)
	b, err := json.MarshalIndent(embedIdQueue, "", "    ")
	if err != nil {
		return err
	}
	var cfg = tool.GetConfig()
	err = ioutil.WriteFile(cfg.Path.Index+strconv.FormatInt(time.Now().Unix(), 10)+".json", b, 0777)
	if err != nil {
		return err
	}
	return nil
}

//将水印图片生成水印切片
func getWatermarkSlice(fileName string) ([]string, error) {
	var cfg = tool.GetConfig()
	imgMat := gocv.IMRead(cfg.Path.Watermark+fileName, gocv.IMReadUnchanged)
	img, err := imgMat.ToImage()
	if err != nil {
		return nil, err
	}
	var watermark = make([]string, 0, imgMat.Rows()*imgMat.Cols())
	for i := 0; i < imgMat.Rows(); i++ {
		for j := 0; j < imgMat.Cols(); j++ {
			r, g, b, a := img.At(i, j).RGBA()
			if (r == g) && (r == b) && (r == a) && (r == uint32(65535)) {
				watermark = append(watermark, "1")
			} else {
				watermark = append(watermark, "0")
			}
		}
	}
	return watermark, nil
}

//得到嵌入索引表结构体
func getIndexStruct(watermark []string, patients []model.Patient) IndexResultSlice {
	var (
		embedIdQueue    IndexResultSlice                // 决定嵌入元组的id号顺序
		pixelCount      = 0                             // 像素点计数
		patientsNums    = len(patients)                 // 总元组个数
		embedPicNums    = patientsNums/20000 + 1        // 需要嵌入的水印图片的个数
		watermarkLength = embedPicNums * len(watermark) // 嵌入像素点总数
		index           = 0                             // 水印嵌入元组的索引值,记录当前水印最后一个像素点的索引,下个水印嵌入从index开始
		temp            IndexResult                     // 临时变量, 最终json文件中的一个结构体
		flag            = 0                             // 用来退出外层循环
		loopTimes       = 0                             // 循环次数,增强代码复用
		existFlag       = 0                             // 判断此元组是否已经被使用过
	)
	rand.Seed(time.Now().Unix())
	for pixelCount < watermarkLength {
		if flag == 1 {
			break
		}
		for i := 0; i < patientsNums; i++ {
			if loopTimes != 0 {
				var tempLength = len(embedIdQueue.Result)
				for k := 0; k < tempLength; k++ {
					if embedIdQueue.Result[k].Id == strconv.FormatInt(patients[i].Id, 10) {
						existFlag = 1
						break
					}
				}
			}
			if existFlag == 1 {
				existFlag = 0
				continue
			} else {
				randNum := rand.Float32()
				if randNum > 0.2 { // 判定条件, 当randNum > 0.2 时就选定该元组
					if res, _ := strconv.Atoi(string(strconv.FormatFloat(rand.Float64(), 'f', 6, 64)[7])); res%2 == 0 {
						temp.Tuple.TupleId = strconv.FormatInt(patients[i].Id, 10)
						temp.Tuple.TupleName = "weight"
						temp.Id = strconv.Itoa(index)
						index += 1
					} else {
						temp.Tuple.TupleId = strconv.FormatInt(patients[i].Id, 10)
						temp.Tuple.TupleName = "bill"
						temp.Id = strconv.Itoa(index)
						index += 1
					}
					embedIdQueue.Result = append(embedIdQueue.Result, temp)
					pixelCount += 1
					if pixelCount == watermarkLength {
						pixelCount = 0
						flag = 1
						break
					}
				} else {
					continue
				}
			}
		}
		loopTimes += 1
	}
	return embedIdQueue
}

//将用户水印嵌入数据库中
func EmbedNewPic(phone string) (int64, error) {
	watermark, err := getWatermarkSlice("secret_new" + phone + ".bmp")
	if err != nil {
		return 0, err
	}
	var patientService service.PatientService
	success, err := patientService.CreateNewUserPatientCopyTable(phone)
	if err != nil {
		return 0, err
	}
	if success {
		indexJson, err := loadJsonIndex()
		if err != nil {
			return 0, err
		}
		patientCopy, err := patientService.GetUserPatientCopyData(phone)
		if err != nil {
			return 0, err
		}
		res, err := embed(watermark, indexJson, patientCopy, phone)
		return res, nil
	} else {
		return 0, err
	}
}

//将原始水印嵌入数据库中
func EmbedOldPic() (int64, error) {
	watermark, err := getWatermarkSlice("secret_old.bmp")
	if err != nil {
		return 0, err
	}
	var patientService service.PatientService
	indexJson, err := loadJsonIndex()
	if err != nil {
		return 0, err
	}
	patientCopy, err := patientService.GetPatientCopyData()
	if err != nil {
		return 0, err
	}
	res, err := embed(watermark, indexJson, patientCopy, "")
	return res, nil
}

//读取索引表Json文件
func loadJsonIndex() (*IndexResultSlice, error) {
	var cfg = tool.GetConfig()
	files, _ := ioutil.ReadDir(cfg.Path.Index)
	f, err := os.Open(cfg.Path.Index + files[0].Name())
	if err != nil {
		return nil, err
	}
	r := io.Reader(f)
	ret := &IndexResultSlice{}
	err = json.NewDecoder(r).Decode(ret)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	return ret, nil
}

//嵌入
func embed(watermark []string, indexJson *IndexResultSlice, patientCopy *[]model.PatientCopy, phone string) (int64,
	error) {
	var length = len(indexJson.Result)
	patients := make(map[string]struct {
		Weight decimal.Decimal
		Bill   decimal.Decimal
	})
	for i := 0; i < len(*patientCopy); i++ {
		patients[strconv.FormatInt((*patientCopy)[i].Id, 10)] = struct {
			Weight decimal.Decimal
			Bill   decimal.Decimal
		}{Weight: (*patientCopy)[i].Weight, Bill: (*patientCopy)[i].Bill}
	}
	for i := 0; i < length; i++ {
		pixel := watermark[i%10000]
		var err error
		if pixel == "0" {
			if indexJson.Result[i].Tuple.TupleName == "weight" {
				if patients[indexJson.Result[i].Tuple.TupleId].Weight.Mul(decimal.NewFromInt(10)).Mod(decimal.
					NewFromInt(2)).String() == "0" {
					continue
				} else {
					var patientCopyService service.PatientService
					id, _ := strconv.ParseInt(indexJson.Result[i].Tuple.TupleId, 10, 64)
					if phone != "" {
						err = patientCopyService.UpdateUserPatientCopy(phone, "weight",
							patients[indexJson.Result[i].Tuple.TupleId].Weight.Sub(decimal.NewFromFloat(0.1)),
							id)
					} else {
						err = patientCopyService.UpdatePatientCopy("weight",
							patients[indexJson.Result[i].Tuple.TupleId].Weight.Sub(decimal.NewFromFloat(0.1)),
							id)
					}
					if err != nil {
						return 0, err
					}
				}
			} else {
				if patients[indexJson.Result[i].Tuple.TupleId].Bill.Mul(decimal.NewFromInt(100)).Mod(decimal.
					NewFromInt(2)).String() == "0" {
					continue
				} else {
					var patientCopyService service.PatientService
					id, _ := strconv.ParseInt(indexJson.Result[i].Tuple.TupleId, 10, 64)
					if phone != "" {
						err = patientCopyService.UpdateUserPatientCopy(phone, "bill",
							patients[indexJson.Result[i].Tuple.TupleId].Bill.Sub(decimal.NewFromFloat(0.01)), id)
					} else {
						err = patientCopyService.UpdatePatientCopy("bill",
							patients[indexJson.Result[i].Tuple.TupleId].Bill.Sub(decimal.NewFromFloat(0.01)), id)
					}
					if err != nil {
						return 0, err
					}
				}
			}
		} else {
			if indexJson.Result[i].Tuple.TupleName == "weight" {
				if patients[indexJson.Result[i].Tuple.TupleId].Weight.Mul(decimal.NewFromInt(10)).Mod(decimal.
					NewFromInt(2)).String() != "0" {
					continue
				} else {
					var patientCopyService service.PatientService
					id, _ := strconv.ParseInt(indexJson.Result[i].Tuple.TupleId, 10, 64)
					if phone != "" {
						err = patientCopyService.UpdateUserPatientCopy(phone, "weight",
							patients[indexJson.Result[i].Tuple.TupleId].Weight.Sub(decimal.NewFromFloat(0.1)), id)
					} else {
						err = patientCopyService.UpdatePatientCopy("weight",
							patients[indexJson.Result[i].Tuple.TupleId].Weight.Sub(decimal.NewFromFloat(0.1)), id)
					}
					if err != nil {
						return 0, err
					}
				}
			} else {
				if patients[indexJson.Result[i].Tuple.TupleId].Bill.Mul(decimal.NewFromInt(100)).Mod(decimal.
					NewFromInt(2)).String() != "0" {
					continue
				} else {
					var patientCopyService service.PatientService
					id, _ := strconv.ParseInt(indexJson.Result[i].Tuple.TupleId, 10, 64)
					if phone != "" {
						err = patientCopyService.UpdateUserPatientCopy(phone, "bill",
							patients[indexJson.Result[i].Tuple.TupleId].Bill.Sub(decimal.NewFromFloat(0.01)), id)
					} else {
						err = patientCopyService.UpdatePatientCopy("bill",
							patients[indexJson.Result[i].Tuple.TupleId].Bill.Sub(decimal.NewFromFloat(0.01)), id)
					}
					if err != nil {
						return 0, err
					}
				}
			}
		}
	}
	if phone != "" {
		var implantIndexService service.ImplantIndexService
		var cfg = tool.GetConfig()
		files, _ := ioutil.ReadDir(cfg.Path.Index)
		iidex := model.ImplantIndex{Phone: phone, IndexName: files[0].Name(), Datetime: strconv.FormatInt(time.Now().Unix(), 10)}
		res, err := implantIndexService.CreateRecord(iidex)
		if err != nil {
			return 0, err
		}
		return res, nil
	} else {
		return 1, nil
	}
}
