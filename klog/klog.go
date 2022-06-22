package klog

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
)

const (
	ColorType_F     = "\x1b[1;31;40m"
	ColorType_A     = "\x1b[91;40m"
	ColorType_C     = "\x1b[1;36;40m"
	ColorType_E     = "\x1b[96;40m"
	ColorType_W     = "\x1b[1;33;40m"
	ColorType_N     = "\x1b[93;40m"
	ColorType_I     = "\x1b[1;32;40m"
	ColorType_D     = "\x1b[92;40m"
	ColorType_Reset = "\x1b[0m"
)

var Conf struct {
	ShortPath bool
	NoColor   bool
	Mute      bool
}

// F :Fatal
func F(formating string, args ...interface{}) {
	color := ColorType_F
	KLogLN(2, Conf.ShortPath, color, "F", formating, args...)
}

// A :Alert
func A(formating string, args ...interface{}) {
	color := ColorType_A
	KLogLN(2, Conf.ShortPath, color, "A", formating, args...)
}

// C :Critical conditions
func C(formating string, args ...interface{}) {
	color := ColorType_C
	KLogLN(2, Conf.ShortPath, color, "C", formating, args...)
}

// E :Error
func E(formating string, args ...interface{}) {
	color := ColorType_E
	KLogLN(2, Conf.ShortPath, color, "E", formating, args...)
}

// W :Warning
func W(formating string, args ...interface{}) {
	color := ColorType_W
	KLogLN(2, Conf.ShortPath, color, "W", formating, args...)
}

// N :Notice
func N(formating string, args ...interface{}) {
	color := ColorType_N
	KLogLN(2, Conf.ShortPath, color, "N", formating, args...)
}

// I :Information
func I(formating string, args ...interface{}) {
	color := ColorType_I
	KLogLN(2, Conf.ShortPath, color, "I", formating, args...)
}

// D :Debug message
func D(formating string, args ...interface{}) {
	color := ColorType_D
	KLogLN(2, Conf.ShortPath, color, "D", formating, args...)
}

// DD :Debug message with depth
func DD(depth int, formating string, args ...interface{}) {
	color := ColorType_D
	KLogLN(depth, Conf.ShortPath, color, "D", formating, args...)
}

func Dump(obj interface{}, strPart ...interface{}) {
	color := ColorType_D

	var s string
	strPartLen := len(strPart)

	switch strPartLen {
	case 0:
		s = spew.Sdump(obj)

	case 1:
		s = strPart[0].(string)
		s += spew.Sdump(obj)

	default:
		fmtPart := strPart[0].(string)
		argPart := strPart[1:len(strPart)]
		s = fmt.Sprintf(fmtPart, argPart...)
		s += spew.Sdump(obj)
	}

	KLog(2, Conf.ShortPath, color, "D", "%s", s)
}

func init() {
	spew.Config.Indent = "    "
	Conf.Mute = false
}
