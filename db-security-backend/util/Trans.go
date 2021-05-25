package util

import (
	"strconv"
	"unsafe"
)

type Trans struct {
}

// Int2Byte Int转Byte
func Int2Byte(data int) (ret []byte) {
	ret = make([]byte, unsafe.Sizeof(data))
	var tmp int = 0xff
	var index uint = 0
	for index = 0; index < uint(unsafe.Sizeof(data)); index++ {
		ret[index] = byte((tmp << (index * 8) & data) >> (index * 8))
	}
	return ret
}

// Byte2Int Byte转Int
func Byte2Int(data []byte) int {
	var ret int = 0
	var i uint = 0
	for i = 0; i < uint(len(data)); i++ {
		ret = ret | (int(data[i]) << (i * 8))
	}
	return ret
}

func PixelToInt(pixel []int) int64 {
	var tmp string
	for p := range pixel {
		if pixel[p] == 255 {
			tmp += "1"
		} else {
			tmp += "0"
		}
	}
	parseInt, err := strconv.ParseInt(tmp, 2, 64)
	if err != nil {
		return 0
	}
	return parseInt
}

// MajorityVoting 多数表决
func MajorityVoting(input []int, threshold int) []int {
	var length = len(input)
	var ans []int
	for i := 0; i < threshold; i++ {
		var count0 = 0
		var count1 = 0
		var j = i
		for j < length {
			if input[i] == 0 {
				count0++
			} else {
				count1++
			}
			j += threshold
		}
		if count1 > count0 {
			ans = append(ans, 255)
		} else {
			ans = append(ans, 0)
		}
	}
	return ans
}
