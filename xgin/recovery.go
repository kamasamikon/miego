package xgin

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http/httputil"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kamasamikon/miego/klog"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func function(pc uintptr) []byte {
	var slash = []byte("/")
	var dot = []byte(".")
	var centerDot = []byte("·")
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}

func source(lines [][]byte, n int) []byte {
	n--
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data

	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

func HandleRecovery(c *gin.Context, err interface{}) {
	var s strings.Builder
	s.WriteString("\n")
	s.WriteString("### PANIC:\n")
	s.WriteString(fmt.Sprintf("%v\n", err))
	s.WriteString("\n")

	httpRequest, _ := httputil.DumpRequest(c.Request, false)
	headers := strings.Split(string(httpRequest), "\r\n")
	s.WriteString("### HEADER:\n")
	for _, header := range headers {
		if header != "" {
			s.WriteString(fmt.Sprintf("1. '%s'\n", header))
		}
	}
	s.WriteString("\n")

	s.WriteString("### STACK:\n")
	s.WriteString("```\n")
	stack := stack(4)
	s.WriteString(fmt.Sprintf("%s\n", stack))
	s.WriteString("```\n")

	klog.E(s.String())

	// SendMessage
}
