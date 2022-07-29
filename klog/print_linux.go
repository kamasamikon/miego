//go:build linux || darwin
// +build linux darwin

package klog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// KLogLN : Log with CR
func KLogLN(dep int, shortPath bool, color string, class string, formating string, args ...interface{}) {
	if Conf.Mute {
		return
	}
	if Conf.NoColor {
		color = ""
	}

	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(dep)

	if ok {
		funcname = runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")
	}

	if shortPath {
		filename = filepath.Base(filename)
	}

	cEnd := ColorType_Reset
	if color == "" {
		cEnd = ""
	}

	now := time.Now()
	nowQ := now.Format("2006/01/02 15:04:05.")
	nowH := now.Nanosecond() / 1000 / 1000 % 1000
	fmt.Printf("%s|%s|S:%s%03d|F:%s|H:%s|L:%d|%s %s\n", color, class, nowQ, nowH, filename, funcname, line, cEnd, fmt.Sprintf(formating, args...))
}

// KLogLN : Log with CR
func KLogLNS(dep int, shortPath bool, color string, class string, formating string, args ...interface{}) string {
	if Conf.Mute {
		return ""
	}
	if Conf.NoColor {
		color = ""
	}

	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(dep)

	if ok {
		funcname = runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")
	}

	if shortPath {
		filename = filepath.Base(filename)
	}

	cEnd := ColorType_Reset
	if color == "" {
		cEnd = ""
	}

	now := time.Now()
	nowQ := now.Format("2006/01/02 15:04:05.")
	nowH := now.Nanosecond() / 1000 / 1000 % 1000
	return fmt.Sprintf("%s|%s|S:%s%03d|F:%s|H:%s|L:%d|%s %s\n", color, class, nowQ, nowH, filename, funcname, line, cEnd, fmt.Sprintf(formating, args...))
}

// KLogX : No '\s' appended.
func KLog(dep int, shortPath bool, color string, class string, formating string, args ...interface{}) {
	if Conf.Mute {
		return
	}
	if Conf.NoColor {
		color = ""
	}

	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(dep)

	if ok {
		funcname = runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")
	}

	if shortPath {
		filename = filepath.Base(filename)
	}

	cEnd := ColorType_Reset
	if color == "" {
		cEnd = ""
	}

	now := time.Now()
	nowQ := now.Format("2006/01/02 15:04:05.")
	nowH := now.Nanosecond() / 1000 / 1000 % 1000
	fmt.Printf("%s|%s|S:%s%03d|F:%s|H:%s|L:%d|%s %s", color, class, nowQ, nowH, filename, funcname, line, cEnd, fmt.Sprintf(formating, args...))
}
