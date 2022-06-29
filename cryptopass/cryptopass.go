package cryptopass

import (
	"crypto/md5"
	"encoding/hex"
	"time"
)

func Make(password string, salt string) string {
	ctx := md5.New()
	ctx.Write([]byte(password))
	ctx.Write([]byte(salt))
	return hex.EncodeToString(ctx.Sum(nil))
}

func Hide(name string) (string, string) {
	Now := time.Now()
	va := Make(name, Now.Format("200601021505"))[0:6]
	vb := Make(name, Now.Format("20060102150405"))

	var output string

	output += va[0:4]
	output += vb[4:30]
	output += va[4:6]

	return va, output
}
