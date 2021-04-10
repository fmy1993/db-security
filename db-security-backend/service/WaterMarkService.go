package service

import (
	"bytes"
	"crypto/sha512"
	"db-security-backend/config"
	"db-security-backend/dao"
	"db-security-backend/model"
	"db-security-backend/util"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/shopspring/decimal"
	"gocv.io/x/gocv"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"strconv"
)

type Img interface {
	Set(x, y int, color color.Color)
}

type WaterMarkService struct {
	ss StaffService
}

// EmbedWatermarkToData 根据数据字典对数据进行更新
func (wms *WaterMarkService) EmbedWatermarkToData(phone string) error {
	dataMap, err := wms.getUpdateDataMap(phone)
	if err != nil {
		return err
	}
	session := config.DbEngine.NewSession()
	defer session.Close()
	_ = session.Begin()
	for key, value := range dataMap {
		sql := fmt.Sprintf("update staff_copy_"+phone+" set salary = %s where staff_id = %d;", value.Salary, key)
		_, _ = session.Exec(sql)
	}
	_ = session.Commit()
	return nil
}

// BlendFingerPrintToPic 将用户指纹嵌入原始水印图片
func (wms *WaterMarkService) BlendFingerPrintToPic(fingerPrint, phone string) {
	cfg := config.GetConfig()
	imgMat := gocv.IMRead(cfg.Path.Watermark+"old_pic.png", gocv.IMReadColor)
	defer imgMat.Close()
	img, err := imgMat.ToImage()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	for i := 0; i < imgMat.Rows(); i++ {
		for j := 0; j < 10; j++ {
			if fingerPrint[j] == '1' {
				img.(Img).Set(j, i, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			} else {
				img.(Img).Set(j, i, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}
	for i := 0; i < imgMat.Rows(); i++ {
		for j := 90; j < 100; j++ {
			if fingerPrint[j-80] == '1' {
				img.(Img).Set(j, i, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			} else {
				img.(Img).Set(j, i, color.RGBA{R: 0, G: 0, B: 0, A: 255})
			}
		}
	}
	demo, err := gocv.ImageToMatRGB(img)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer demo.Close()
	gocv.IMWrite(cfg.Path.Watermark+"new_pic"+phone+".bmp", demo)
}

// PickWatermarkByCsv 根据csv数据提取水印
func (wms *WaterMarkService) PickWatermarkByCsv(staffCopys map[int64]decimal.Decimal) {
	var cfg = config.GetConfig()
	tuples, _ := wms.chooseTupleById()
	var tuplesLength = len(*tuples)
	var wm []int
	for i := 0; i < tuplesLength; i++ {
		tmp, exist := staffCopys[(*tuples)[i]]
		if !exist {
			wm = append(wm, 255)
			continue
		}
		tmp = staffCopys[(*tuples)[i]].Mod(decimal.NewFromInt(2))
		if tmp.Equal(decimal.NewFromInt(1)) {
			wm = append(wm, 255)
		} else {
			wm = append(wm, 0)
		}
	}
	var wmNums = len(wm) / 10000
	for k := 0; k < wmNums; k++ {
		var imgMat = gocv.NewMatWithSize(100, 100, gocv.MatTypeCV8UC3)
		img, _ := imgMat.ToImage()
		for i := 0; i < 100; i++ {
			for j := 0; j < 100; j++ {
				if (wm)[i*100+j+k*10000] == 0 {
					img.(Img).Set(i, j, color.RGBA{R: 0, G: 0, B: 0, A: 255})
				} else {
					img.(Img).Set(i, j, color.RGBA{R: 255, G: 255, B: 255, A: 255})
				}
			}
		}
		demo, _ := gocv.ImageToMatRGB(img)
		gocv.IMWrite(cfg.Path.PickUp+"res"+strconv.Itoa(k)+".bmp", demo)
		imgMat.Close()
		demo.Close()
	}
	wms.DeArnold(3)
}

// PickFingerPrintByPic 根据提取到的水印图片提取指纹
func (wms *WaterMarkService) PickFingerPrintByPic() string {
	var cfg = config.GetConfig()
	fileList, _ := ioutil.ReadDir(cfg.Path.Res)
	var fp [][]string
	for i := range fileList {
		imgMat := gocv.IMRead(cfg.Path.Res+fileList[i].Name(), gocv.IMReadUnchanged)
		img, _ := imgMat.ToImage()
		for i := 0; i < imgMat.Rows(); i++ {
			var temp []string
			for j := 0; j < 10; j++ {
				r, g, b, _ := img.At(j, i).RGBA()
				if r == 65535 && g == 65535 && b == 65535 {
					temp = append(temp, "1")
				} else {
					temp = append(temp, "0")
				}
			}
			for k := 90; k < 100; k++ {
				r, g, b, _ := img.At(k, i).RGBA()
				if r == 65535 && g == 65535 && b == 65535 {
					temp = append(temp, "1")
				} else {
					temp = append(temp, "0")
				}
			}
			fp = append(fp, temp)
		}
	}
	var finalFp string
	for i := 0; i < 20; i++ {
		var count0 = 0
		var count1 = 0
		for j := 0; j < len(fp); j++ {
			if (fp)[j][i] == "1" {
				count1 += 1
			} else {
				count0 += 1
			}
		}
		if count1 > count0 {
			finalFp += "1"
		} else {
			finalFp += "0"
		}
	}
	return finalFp
}

// Arnold Arnold变换
func (wms *WaterMarkService) Arnold(phone string, key int) {
	cfg := config.GetConfig()
	var imgMat gocv.Mat
	if phone == "" {
		imgMat = gocv.IMRead(cfg.Path.Watermark+"old_pic.png", gocv.IMReadUnchanged)
	} else {
		imgMat = gocv.IMRead(cfg.Path.Watermark+"new_pic"+phone+".bmp", gocv.IMReadUnchanged)
	}
	defer imgMat.Close()
	img, err := imgMat.ToImage()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	tempMat := gocv.NewMatWithSize(imgMat.Rows(), imgMat.Cols(), gocv.MatTypeCV8UC3)
	tempImg, _ := tempMat.ToImage()
	for k := 0; k < key; k++ {
		for i := 0; i < imgMat.Rows(); i++ {
			for j := 0; j < imgMat.Cols(); j++ {
				tempImg.(Img).Set((i+j)%imgMat.Rows(), (i+2*j)%imgMat.Rows(), img.At(i, j))
			}
		}
		for i := 0; i < imgMat.Rows(); i++ {
			for j := 0; j < imgMat.Cols(); j++ {
				img.(Img).Set(i, j, tempImg.At(i, j))
			}
		}
	}
	demo, err := gocv.ImageToMatRGB(img)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer demo.Close()
	if phone == "" {
		gocv.IMWrite(cfg.Path.Watermark+"secret_old.bmp", demo)
	} else {
		gocv.IMWrite(cfg.Path.Watermark+"secret_new"+phone+".bmp", demo)
	}
}

// DeArnold 逆Arnold变换
func (wms *WaterMarkService) DeArnold(key int) {
	var cfg = config.GetConfig()
	files, _ := ioutil.ReadDir(cfg.Path.PickUp)
	for count := 0; count < len(files); count++ {
		imgMat := gocv.IMRead(cfg.Path.PickUp+"res"+strconv.Itoa(count)+".bmp", gocv.IMReadUnchanged)
		img, err := imgMat.ToImage()
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		tempMat := gocv.NewMatWithSize(imgMat.Rows(), imgMat.Cols(), gocv.MatTypeCV8UC3)
		tempImg, _ := tempMat.ToImage()
		for k := 0; k < key; k++ {
			for i := 0; i < imgMat.Rows(); i++ {
				for j := 0; j < imgMat.Cols(); j++ {
					var tempI, tempJ int
					if 2*i-j < 0 {
						tempI = imgMat.Rows() + (2*i - j)
					} else {
						tempI = (2*i - j) % imgMat.Rows()
					}
					if j-i < 0 {
						tempJ = imgMat.Rows() + j - i
					} else {
						tempJ = (j - i) % imgMat.Rows()
					}
					tempImg.(Img).Set(tempI, tempJ, img.At(i, j))
				}
			}
			for i := 0; i < imgMat.Rows(); i++ {
				for j := 0; j < imgMat.Cols(); j++ {
					img.(Img).Set(i, j, tempImg.At(i, j))
				}
			}
		}
		demo, err := gocv.ImageToMatRGB(img)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
		gocv.IMWrite(cfg.Path.Res+"res_pic_"+strconv.Itoa(count)+".bmp", demo)
		imgMat.Close()
		demo.Close()
	}
}

// PicDecode 水印图像编码
func (wms *WaterMarkService) PicDecode() string {
	var cfg = config.GetConfig()
	ffMat := gocv.IMRead(cfg.Path.Res+"res_pic_0.bmp", gocv.IMReadUnchanged)
	ff, _ := ffMat.ToImage()
	emptyBuff := bytes.NewBuffer(nil)
	_ = jpeg.Encode(emptyBuff, ff, nil)
	dist := make([]byte, 50000)
	base64.StdEncoding.Encode(dist, emptyBuff.Bytes())
	index := bytes.IndexByte(dist, 0)
	baseImage := dist[0:index]
	return "data:image/png;base64," + string(baseImage)
}

// 选择需要嵌入的元祖
func (wms *WaterMarkService) chooseTupleById() (*[]int64, error) {
	var staffs *[]model.Staff
	staffs, err := wms.ss.GetAllStaff()
	var res []int64
	if err != nil {
		return nil, err
	}
	var n = len(*staffs)
	var priKey = wms.getPrimaryKey()
	for i := 0; i < n; i++ {
		var a, _ = hex.DecodeString(fmt.Sprintf("%x", sha512.Sum512([]byte(strconv.Itoa(int((*staffs)[i].StaffId^priKey))))))
		var c = sha512.Sum512([]byte(strconv.Itoa(util.Byte2Int(a) ^ int(priKey))))
		var d, _ = hex.DecodeString(fmt.Sprintf("%x", c))
		var e = util.Byte2Int(d)
		if e%2 == 0 {
			res = append(res, (*staffs)[i].StaffId)
		}
	}
	return &res, nil
}

// 获取密钥
func (wms *WaterMarkService) getPrimaryKey() int64 {
	var configDao = dao.NewConfigDao()
	cfg, _ := configDao.GetConfigByColumn("pri_key")
	res, _ := strconv.Atoi(cfg.ConfigValue)
	return int64(res)
}

// 将水印图片生成水印切片
func (wms *WaterMarkService) getWatermarkSlice(fileName string) ([]string, error) {
	var cfg = config.GetConfig()
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

// 得到更新后的用户表数据字典
func (wms *WaterMarkService) getUpdateDataMap(phone string) (map[int64]model.StaffCopy, error) {
	watermark, err := wms.getWatermarkSlice("secret_new" + phone + ".bmp")
	if err != nil {
		return nil, err
	}
	var staffCopyService StaffCopyService
	staffCopyService.IsUserStaffCopyExist(phone)
	staffCopys, _ := staffCopyService.GetAllStaffCopyData()
	var staffCopysLength = len(*staffCopys)
	var staffCopysMap = make(map[int64]model.StaffCopy, staffCopysLength)
	for i := 0; i < staffCopysLength; i++ {
		staffCopysMap[(*staffCopys)[i].StaffId] = (*staffCopys)[i]
	}
	var tuples *[]int64
	tuples, _ = wms.chooseTupleById()
	var tupPtr int64 = 0
	var wmNums = len(*tuples) / 10000
	for i := 0; i < wmNums; i++ {
		var waterMarkPtr = 0
		for ; waterMarkPtr < len(watermark); waterMarkPtr++ {
			if watermark[waterMarkPtr] == "1" {
				if staffCopysMap[(*tuples)[tupPtr]].Salary.Mod(decimal.NewFromInt(2)).Equal(decimal.NewFromInt(0)) {
					tmp := staffCopysMap[(*tuples)[tupPtr]]
					tmp.Salary = tmp.Salary.Add(decimal.NewFromInt(1))
					staffCopysMap[(*tuples)[tupPtr]] = tmp
				}
			} else {
				if staffCopysMap[(*tuples)[tupPtr]].Salary.Mod(decimal.NewFromInt(2)).Equal(decimal.NewFromInt(1)) {
					tmp := staffCopysMap[(*tuples)[tupPtr]]
					tmp.Salary = tmp.Salary.Add(decimal.NewFromInt(1))
					staffCopysMap[(*tuples)[tupPtr]] = tmp
				}
			}
			tupPtr++
		}
	}
	return staffCopysMap, nil
}
