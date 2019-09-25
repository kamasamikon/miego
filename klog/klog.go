package klog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mgutz/ansi"
)

var crF = ansi.ColorCode("red+b:black")
var crA = ansi.ColorCode("red+h:black")
var crC = ansi.ColorCode("cyan+b:black")
var crE = ansi.ColorCode("cyan+h:black")
var crW = ansi.ColorCode("yellow+b:black")
var crN = ansi.ColorCode("yellow+h:black")
var crI = ansi.ColorCode("green+b:black")
var crD = ansi.ColorCode("green+h:black")
var reset = ansi.ColorCode("reset")

var ctrlMap map[string]bool

var Conf struct {
	ShortPath bool
	NoClor    bool
}

func klog(color string, class string, formating string, args ...interface{}) {

	filename, line, funcname := "???", 0, "???"
	pc, filename, line, ok := runtime.Caller(2)

	if ok {
		funcname = runtime.FuncForPC(pc).Name()
		funcname = filepath.Ext(funcname)
		funcname = strings.TrimPrefix(funcname, ".")
	}

	if Conf.ShortPath {
		filename = filepath.Base(filename)
	}
	if Conf.NoClor {
		color = ""
	}

	/*
		{
			keyA := fmt.Sprintf("%d@%s", line, funcname)
			keyB := fmt.Sprintf("%d@%s", "*", funcname)
			keyC := fmt.Sprintf("%d@%s", line, "*")
			keyD := fmt.Sprintf("%d@%s", "*", "*")
		}
	*/

	now := time.Now()
	nowQ := now.Format("2006/01/02 15:04:05.")
	nowH := now.Nanosecond() / 1000 / 1000 % 1000
	fmt.Printf("%s|%s|S:%s%03d|F:%s|H:%s|L:%d|%s %s\n", color, class, nowQ, nowH, filename, funcname, line, reset, fmt.Sprintf(formating, args...))
}

// F :Fatal
func F(formating string, args ...interface{}) {
	klog(crF, "F", formating, args...)
}

// A :Alert
func A(formating string, args ...interface{}) {
	klog(crA, "A", formating, args...)
}

// C :Critical conditions
func C(formating string, args ...interface{}) {
	klog(crC, "C", formating, args...)
}

// E :Error
func E(formating string, args ...interface{}) {
	klog(crE, "E", formating, args...)
}

// W :Warning
func W(formating string, args ...interface{}) {
	klog(crW, "W", formating, args...)
}

// N :Notice
func N(formating string, args ...interface{}) {
	klog(crN, "N", formating, args...)
}

// I :Information
func I(formating string, args ...interface{}) {
	klog(crI, "I", formating, args...)
}

// D :Debug message
func D(formating string, args ...interface{}) {
	klog(crD, "D", formating, args...)
}
