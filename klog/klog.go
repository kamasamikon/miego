package klog

import (
	"fmt"
	"io"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

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

/////////////////////////////////////////////////////////////////////////

var Conf struct {
	// 0:false, 1:true, -1:notset
	ShortPath int // 只显示文件名，否则会显示文件的全路径
	Dull      int // 单色的，不是彩色的
	Mute      int // 是否把文字真的打印出来
	UseStdout int // WriterAdd(Stdout)
	Writers   map[string]io.Writer
}

///////////////////////////////////////////////////////////////////////////

// XXX 没有加保护，只在启动的时候弄一下
func WriterAdd(name string, writer io.Writer) {
	Conf.Writers[name] = writer
}

// XXX 没有加保护
func WriterGet(name string) io.Writer {
	w, _ := Conf.Writers[name]
	return w
}

// XXX 没有加保护
func WriterRem(name string) {
	delete(Conf.Writers, name)
}

///////////////////////////////////////////////////////////////////////////

// KLogS : KLog as LN (Line) to string
// @lf: append line feed
func KLogS(dep int, shortPath int, color string, class rune, lf bool, formating string, args ...interface{}) string {
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
	sb.WriteRune(class)

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
func KLogLN(dep int, shortPath int, color string, class rune, formating string, args ...interface{}) {
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
func KLog(dep int, shortPath int, color string, class rune, formating string, args ...interface{}) {
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

/////////////////////////////////////////////////////////////////////////

// S :String
func S(formating string, args ...interface{}) string {
	return KLogS(2, Conf.ShortPath, "", 'S', true, formating, args...)
}

// S :String with Color
func SC(color string, formating string, args ...interface{}) string {
	return KLogS(2, Conf.ShortPath, color, 'S', true, formating, args...)
}

////////////////////////////////////////////////////////////////////////////

// F :Fatal
func F(formating string, args ...interface{}) {
	color := ColorType_F
	KLogLN(2, Conf.ShortPath, color, 'F', formating, args...)
}

// A :Alert
func A(formating string, args ...interface{}) {
	color := ColorType_A
	KLogLN(2, Conf.ShortPath, color, 'A', formating, args...)
}

// C :Critical conditions
func C(formating string, args ...interface{}) {
	color := ColorType_C
	KLogLN(2, Conf.ShortPath, color, 'C', formating, args...)
}

// E :Error
func E(formating string, args ...interface{}) {
	color := ColorType_E
	KLogLN(2, Conf.ShortPath, color, 'E', formating, args...)
}

// W :Warning
func W(formating string, args ...interface{}) {
	color := ColorType_W
	KLogLN(2, Conf.ShortPath, color, 'W', formating, args...)
}

// N :Notice
func N(formating string, args ...interface{}) {
	color := ColorType_N
	KLogLN(2, Conf.ShortPath, color, 'N', formating, args...)
}

// I :Information
func I(formating string, args ...interface{}) {
	color := ColorType_I
	KLogLN(2, Conf.ShortPath, color, 'I', formating, args...)
}

// D :Debug message
func D(formating string, args ...interface{}) {
	color := ColorType_D
	KLogLN(2, Conf.ShortPath, color, 'D', formating, args...)
}

// DD :Debug message with depth
func DD(depth int, formating string, args ...interface{}) {
	color := ColorType_D
	KLogLN(depth, Conf.ShortPath, color, 'D', formating, args...)
}

////////////////////////////////////////////////////////////////////////////

func Color(color string, args ...interface{}) string {
	if len(args) > 0 {
		format := args[0].(string)
		return "\033[" + color + "m" + fmt.Sprintf(format, args[1:]...) + "\033[0m"
	} else {
		return "\033[" + color + "m"
	}
}

///////////////////////////////////////////////////////////////////////////

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
		argPart := strPart[1:]
		s = fmt.Sprintf(fmtPart, argPart...)
		s += cfg.Sdump(obj)
	}

	return s
}

func Dump(obj interface{}, strPart ...interface{}) {
	if Conf.Mute == 1 {
		return
	}

	color := ColorType_D
	s := DumpS(obj, strPart...)
	KLog(2, Conf.ShortPath, color, 'D', "%s", s)
}
