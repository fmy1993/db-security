package util

import "unsafe"

type Trans struct {
}

// Int转Byte
func Int2Byte(data int) (ret []byte) {
	ret = make([]byte, unsafe.Sizeof(data))
	var tmp int = 0xff
	var index uint = 0
	for index = 0; index < uint(unsafe.Sizeof(data)); index++ {
		ret[index] = byte((tmp << (index * 8) & data) >> (index * 8))
	}
	return ret
}

// Byte转Int
func Byte2Int(data []byte) int {
	var ret int = 0
	var i uint = 0
	for i = 0; i < uint(len(data)); i++ {
		ret = ret | (int(data[i]) << (i * 8))
	}
	return ret
}
