package klog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// BT : Print back trace
func BT(maxdep int, formating string, args ...interface{}) {
	if maxdep < 2 {
		KLogLN(2, Conf.ShortPath, ColorType_N, 'T', formating, args...)
		return
	}
	if Conf.Mute == 1 {
		return
	}
	if Conf.Dull == 1 {
		return
	}

	cEnd := ColorType_Reset
	cStart := ColorType_N

	txt := fmt.Sprintf(formating, args...)
	content := fmt.Sprintf("%s|BT|%s %s\n", cStart, cEnd, txt)

	dep := 0
	for {
		dep++

		if maxdep > 0 && dep > maxdep {
			break
		}

		pc, filename, line, ok := runtime.Caller(dep)
		if !ok {
			break
		}

		funcname := runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")

		content += fmt.Sprintf("%s|%02d| F:%s|H:%s|L:%d|%s\n", cStart, dep, filename, funcname, line, cEnd)
	}

	s := []byte(content)
	Writers := Conf.Writers
	for _, w := range Writers {
		w.Write(s)
	}
}
