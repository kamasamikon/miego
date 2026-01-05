package klog

import (
	"path/filepath"
	"runtime"
	"strings"
)

// 渐进式一条条覆盖前边的规则
type FilterRule struct {
	class    rune   // facewind
	fileName string // fileName="^mbs.go"
	funcName string // funcName="^mbs.go"
	line     int    // line=33
	enable   bool   // 是否打印
}

var filterRules []*FilterRule

// 这里的问题是每次都计算，实际上不必每次都算
// 规则发生编号的时候计算就可以了
//
// 目前版本先都计算，主要看这个玩意儿是否可以
func Enabled(fileName, funcName string, class rune, line int) bool {
	// 初始化的规则：fileName = "", funcName = "", line = -1, enable = true
	enable := false
	for _, r := range filterRules {
		enable = false
		if r.class != '*' && r.class != class {
			continue
		}
		if r.fileName != "" && r.fileName != fileName {
			continue
		}
		if r.funcName != "" && r.funcName != funcName {
			continue
		}
		if r.line != -1 && r.line != line {
			continue
		}
		enable = r.enable
	}
	return enable
}

func xKLogS(dep int, shortPath int, color string, class rune, lf bool, formating string, args ...interface{}) string {
	pc, filename, line, _ := runtime.Caller(dep)
	funcname := runtime.FuncForPC(pc).Name()
	funcname = filepath.Ext(funcname)
	funcname = strings.TrimPrefix(funcname, ".")

	if !Enabled(filename, funcname, class, line) {
		return ""
	}
	return ""
}

func yKLogS(fcache *int, dep int, shortPath int, color string, class rune, lf bool, formating string, args ...interface{}) string {
	// if *fcache == len(filterRules) {
	// Enabled := *Enabled
	// }

	pc, filename, line, _ := runtime.Caller(dep)
	funcname := runtime.FuncForPC(pc).Name()
	funcname = filepath.Ext(funcname)
	funcname = strings.TrimPrefix(funcname, ".")

	if !Enabled(filename, funcname, class, line) {
		return ""
	}

	// ... go to print
	return ""
}
