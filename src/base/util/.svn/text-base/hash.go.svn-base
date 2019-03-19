package util

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"strings"
)

func BKDRHash(str string) uint32 {

	str = strings.ToUpper(str)

	var seed uint32 = 131 // the magic number, 31, 131, 1313, 13131, etc.. orz..
	var hash uint32 = 0

	for i := 0; i != len(str); i++ {
		hash = hash*seed + uint32(str[i])
	}

	return hash
}

func Md5(str string) string {
	md5 := md5.New()
	io.WriteString(md5, str)
	return hex.EncodeToString(md5.Sum(nil))
}
