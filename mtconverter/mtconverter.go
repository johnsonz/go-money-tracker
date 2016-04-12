package mtconverter

import "strconv"

func Bytes2Float64(b []byte) (float64, error) {
	return strconv.ParseFloat(string(b), 64)
}
func Float642String(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}
func Int642String(i int64) string {
	return strconv.FormatInt(i, 10)
}
func Bytes2Int(b []byte) (int64, error) {
	return strconv.ParseInt(string(b), 10, 64)
}
