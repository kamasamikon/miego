package cryptopass

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

func Make(password string, salt string) string {
	fmt.Printf("%s %s\n", password, salt)
	ctx := md5.New()
	ctx.Write([]byte(password))
	ctx.Write([]byte(salt))
	return hex.EncodeToString(ctx.Sum(nil))
}

func Hide(name string) (string, string) {
	Now := time.Now().Format("20060102150405")
	va := Make(name, Now[0:12])[0:6]
	vb := Make(name, Now)

	var output string

	output += va[0:4]
	output += vb[4:30]
	output += va[4:6]

	return va, output
}
