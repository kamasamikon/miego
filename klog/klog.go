package klog

import (
	"fmt"
	"io"
	"os"

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
	Writers   []io.Writer
}

// XXX 没有加保护
func WriterAdd(args ...interface{}) {
	for _, arg := range args {
		if w, ok := arg.(string); ok {
			switch w {
			case "stdout":
				Conf.Writers = append(Conf.Writers, os.Stdout)
			case "stderr":
				Conf.Writers = append(Conf.Writers, os.Stderr)
			}
			continue
		}

		if w, ok := arg.(io.Writer); ok {
			Conf.Writers = append(Conf.Writers, w)
			continue
		}
	}
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

func Color(color string, formating string, args ...interface{}) string {
	return "\033[0;" + color + "m" + fmt.Sprintf(formating, args...) + "\033[0m"
}

func DumpS(obj interface{}, strPart ...interface{}) string {
	cfg := spew.ConfigState{SortKeys: true, Indent: "    "}

	var s string
	strPartLen := len(strPart)

	switch strPartLen {
	case 0:
		s = cfg.Sdump(obj)

	case 1:
		s = strPart[0].(string)
		s += cfg.Sdump(obj)

	default:
		fmtPart := strPart[0].(string)
		argPart := strPart[1:len(strPart)]
		s = fmt.Sprintf(fmtPart, argPart...)
		s += cfg.Sdump(obj)
	}

	return s
}

func Dump(obj interface{}, strPart ...interface{}) {
	color := ColorType_D
	s := DumpS(obj, strPart...)
	KLog(2, Conf.ShortPath, color, "D", "%s", s)
}

func init() {
	Conf.Mute = false

	WriterAdd("stdout")
}
