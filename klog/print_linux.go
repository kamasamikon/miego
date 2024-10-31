//go:build linux

package klog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// KLogLN : Log with CR
func KLogLN(dep int, shortPath int, color string, class string, formating string, args ...interface{}) {
	if Conf.Mute == 1 {
		return
	}
	if Conf.Dull == 1 {
		color = ""
	}

	pc, filename, line, _ := runtime.Caller(dep)

	funcname := runtime.FuncForPC(pc).Name()
	funcname = filepath.Ext(funcname)
	funcname = strings.TrimPrefix(funcname, ".")

	if shortPath == 1 {
		filename = filepath.Base(filename)
	}

	cEnd := ColorType_Reset
	if color == "" {
		cEnd = ""
	}

	now := time.Now().Format("2006/01/02 15:04:05.000000")
	s := []byte(fmt.Sprintf("%s|%s|S:%s|F:%s|H:%s|L:%d|%s %s\n", color, class, now, filename, funcname, line, cEnd, fmt.Sprintf(formating, args...)))

	Writers := Conf.Writers
	for _, w := range Writers {
		w.Write(s)
	}
}

// KLogLN : Log with CR
func KLogLNS(dep int, shortPath bool, color string, class string, formating string, args ...interface{}) string {
	if Conf.Mute == 1 {
		return ""
	}
	if Conf.Dull == 1 {
		color = ""
	}

	pc, filename, line, _ := runtime.Caller(dep)

	funcname := runtime.FuncForPC(pc).Name()
	funcname = filepath.Ext(funcname)
	funcname = strings.TrimPrefix(funcname, ".")

	if shortPath {
		filename = filepath.Base(filename)
	}

	cEnd := ColorType_Reset
	if color == "" {
		cEnd = ""
	}

	now := time.Now().Format("2006/01/02 15:04:05.000000")
	return fmt.Sprintf("%s|%s|S:%s|F:%s|H:%s|L:%d|%s %s\n", color, class, now, filename, funcname, line, cEnd, fmt.Sprintf(formating, args...))
}

// KLogX : No '\s' appended.
func KLog(dep int, shortPath int, color string, class string, formating string, args ...interface{}) {
	if Conf.Mute == 1 {
		return
	}
	if Conf.Dull == 1 {
		color = ""
	}

	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(dep)

	if ok {
		funcname = runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")
	}

	if shortPath == 1 {
		filename = filepath.Base(filename)
	}

	cEnd := ColorType_Reset
	if color == "" {
		cEnd = ""
	}

	now := time.Now().Format("2006/01/02 15:04:05.000000")
	s := []byte(fmt.Sprintf("%s|%s|S:%s|F:%s|H:%s|L:%d|%s %s", color, class, now, filename, funcname, line, cEnd, fmt.Sprintf(formating, args...)))

	Writers := Conf.Writers
	for _, w := range Writers {
		w.Write(s)
	}
}
