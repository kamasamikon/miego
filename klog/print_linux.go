//go:build linux

package klog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// KLogS : KLog as LN (Line) to string
// @lf: append line feed
func KLogS(dep int, shortPath int, color string, class string, lf bool, formating string, args ...interface{}) string {
	pc, filename, line, _ := runtime.Caller(dep)

	funcname := runtime.FuncForPC(pc).Name()
	funcname = filepath.Ext(funcname)
	funcname = strings.TrimPrefix(funcname, ".")

	if shortPath == 1 {
		filename = filepath.Base(filename)
	}

	now := time.Now().Format("2006/01/02 15:04:05.000000")

	var sb strings.Builder

	sb.WriteString(color)

	sb.WriteRune('|')
	sb.WriteString(class)

	sb.WriteRune('|')
	sb.WriteString(now)

	sb.WriteRune('|')
	sb.WriteString(filename)

	sb.WriteRune('|')
	sb.WriteString(funcname)

	sb.WriteRune('|')
	sb.WriteString(strconv.Itoa(line))

	sb.WriteRune('|')

	if color != "" {
		sb.WriteString(ColorType_Reset)
	}

	sb.WriteRune(' ')
	sb.WriteString(fmt.Sprintf(formating, args...))

	if lf {
		sb.WriteString("\n")
	}

	return sb.String()
}

// KLogLN : KLog as LN (Line)
func KLogLN(dep int, shortPath int, color string, class string, formating string, args ...interface{}) {
	if Conf.Mute == 1 {
		return
	}
	if Conf.Dull == 1 {
		return
	}

	s := []byte(KLogS(dep+1, shortPath, color, class, true, formating, args...))

	Writers := Conf.Writers
	for _, w := range Writers {
		w.Write(s)
	}
}

// KLogX : No '\s' appended.
func KLog(dep int, shortPath int, color string, class string, formating string, args ...interface{}) {
	if Conf.Mute == 1 {
		return
	}
	if Conf.Dull == 1 {
		color = ""
	}

	s := []byte(KLogS(dep+1, shortPath, color, class, false, formating, args...))

	Writers := Conf.Writers
	for _, w := range Writers {
		w.Write(s)
	}
}
