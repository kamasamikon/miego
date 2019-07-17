package cryptopass

import (
	"crypto/md5"
	"encoding/hex"
)

func Make(password string, salt string) string {
	ctx := md5.New()
	ctx.Write([]byte(password))
	ctx.Write([]byte(salt))
	return hex.EncodeToString(ctx.Sum(nil))
}
