package mtconverter

import (
	"bytes"
	"encoding/binary"
	"math"
	"strconv"
)

func Bytes2Float64(b []byte) float64 {
	bits := binary.LittleEndian.Uint64(b)
	float := math.Float64frombits(bits)
	return float
}

func Float642Bytes(float float64) []byte {
	bits := math.Float64bits(float)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, bits)
	return b
}
func Int2Bytes(i int) []byte {
	return []byte(strconv.Itoa(i))
}
func Bytes2Int(b []byte) (int, error) {
	return strconv.Atoi(string(b[:bytes.Index(b, []byte{0})]))
}
