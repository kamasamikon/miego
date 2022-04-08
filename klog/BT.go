package klog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// BT : Print back trace
func BT(maxdep int, formating string, args ...interface{}) string {
	if Conf.Mute {
		return ""
	}

	now := time.Now()
	nowQ := now.Format("2006/01/02 15:04:05.")
	nowH := now.Nanosecond() / 1000 / 1000 % 1000

	cEnd := ColorType_Reset
	cStart := ColorType_D

	content := ""

	txt := fmt.Sprintf(formating, args...)
	content += fmt.Sprintf("%s|BT>>>>|%s %s\n", cStart, cEnd, txt)

	dep := 0
	for {
		dep += 1

		if maxdep > 0 && dep > maxdep {
			break
		}

		pc, filename, line, ok := runtime.Caller(dep)
		if ok == false {
			content += fmt.Sprintf("%s|BT<<<<|%s %s\n", cStart, cEnd, txt)
			return content
		}

		funcname := runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")

		content += fmt.Sprintf("%s|BT|S:%s%03d|F:%s|H:%s|L:%d|%s\n", cStart, nowQ, nowH, filename, funcname, line, cEnd)
	}

	content += fmt.Sprintf("%s|BT<<<<|%s %s\n", cStart, cEnd, txt)
	return content
}
