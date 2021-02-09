package util

import (
	"bytes"
	"encoding/base64"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"log"
	"strconv"

	"db-security-backend/tool"
	"github.com/shopspring/decimal"
	"gocv.io/x/gocv"
)

type Img interface {
	Set(x, y int, color color.Color)
}

type Picture struct {
}

//将指纹嵌入原始图片中
func (p *Picture) BlendFP(fingerPrint, phone string) {
	cfg := tool.GetConfig()
	imgMat := gocv.IMRead(cfg.Path.Watermark+"old_pic.bmp", gocv.IMReadColor)
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

//Arnold变换
func (p *Picture) Arnold(phone string, key int) {
	cfg := tool.GetConfig()
	var imgMat gocv.Mat
	if phone == "" {
		imgMat = gocv.IMRead(cfg.Path.Watermark+"old_pic.bmp", gocv.IMReadUnchanged)
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

//逆Arnold变换
func (p *Picture) DeArnold(key int) {
	var cfg = tool.GetConfig()
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

//提取水印
func (p *Picture) PickWatermark(patientCopy map[string]struct {
	Weight decimal.Decimal
	Bill   decimal.Decimal
}) string {
	var cfg = tool.GetConfig()
	indexJson, _ := loadJsonIndex()
	var wm []int
	length := len(indexJson.Result)
	for i := 0; i < length; i++ {
		if indexJson.Result[i].Tuple.TupleName == "weight" {
			if patientCopy[indexJson.Result[i].Tuple.TupleId].Weight.Mul(decimal.NewFromInt(10)).Mod(decimal.
				NewFromInt(2)).String() == "0" {
				wm = append(wm, 0)
			} else {
				wm = append(wm, 255)
			}
		} else {
			if patientCopy[indexJson.Result[i].Tuple.TupleId].Bill.Mul(decimal.NewFromInt(100)).Mod(decimal.
				NewFromInt(2)).String() == "0" {
				wm = append(wm, 0)
			} else {
				wm = append(wm, 255)
			}
		}
	}
	var finalWm []int
	for i := 10000; i < 15000; i++ {
		finalWm = append(finalWm, (wm)[i])
	}
	for i := 5000; i < 10000; i++ {
		finalWm = append(finalWm, (wm)[i])
	}
	for i := 0; i < 10000; i++ {
		finalWm = append(finalWm, finalWm[i])
	}
	for k := 0; k < len(finalWm)/10000; k++ {
		var imgMat = gocv.NewMatWithSize(100, 100, gocv.MatTypeCV8UC3)
		img, _ := imgMat.ToImage()
		for i := 0; i < 100; i++ {
			for j := 0; j < 100; j++ {
				if (finalWm)[i*100+j+k*10000] == 0 {
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
	p.DeArnold(3)
	return p.PickFp()
}

//提取指纹
func (p *Picture) PickFp() string {
	var cfg = tool.GetConfig()
	imgMat := gocv.IMRead(cfg.Path.Res+"res_pic_0.bmp", gocv.IMReadUnchanged)
	img, _ := imgMat.ToImage()
	var fp [][]string
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

func (p *Picture) PicDecode() string {
	var cfg = tool.GetConfig()
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
