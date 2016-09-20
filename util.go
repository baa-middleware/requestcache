package requestcache

import (
	"crypto/md5"
	"encoding/hex"
)

func md5Encode(str string) string {
	hexStr := md5.Sum([]byte(str))
	return hex.EncodeToString(hexStr[:])
}
