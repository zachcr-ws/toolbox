package code

import (
	"crypto/md5"
	"encoding/hex"
)

func Md5Encode(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}
