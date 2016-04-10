package mtcrypto

import "crypto/md5"

func MD5(input string) [16]byte {
	data := []byte(input)
	return md5.Sum(data)
}
