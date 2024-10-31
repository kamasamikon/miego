package klog

import (
	"io"
	"os"
	"strings"

	"miego/conf"
)

func init() {
	Conf.ShortPath = -1
	Conf.Dull = -1
	Conf.Mute = -1
	Conf.UseStdout = -1
	Conf.Writers = make(map[string]io.Writer)

	// 优先级: 命令行 > 环境变量 > 配置文件
	LoadFromConf()
	LoadFromEnv()
	LoadFromArg()

	if Conf.ShortPath == -1 {
		Conf.ShortPath = 0
	}
	if Conf.Dull == -1 {
		Conf.Dull = 0
	}
	if Conf.Mute == -1 {
		Conf.Mute = 0
	}
	if Conf.UseStdout == -1 {
		Conf.UseStdout = 1
		WriterAdd("stdout", os.Stdout)
	}
}

func LoadFromConf() {
	if x, ok := conf.BoolX("b:/klog/shortPath"); ok {
		if x {
			Conf.ShortPath = 1
		} else {
			Conf.ShortPath = 0
		}
	}

	if x, ok := conf.BoolX("b:/klog/dull"); ok {
		if x {
			Conf.Dull = 1
		} else {
			Conf.Dull = 0
		}
	}

	if x, ok := conf.BoolX("b:/klog/mute"); ok {
		if x {
			Conf.Mute = 1
		} else {
			Conf.Mute = 0
		}
	}

	if x, ok := conf.BoolX("b:/klog/useStdout"); ok {
		if x {
			Conf.UseStdout = 1
		} else {
			Conf.UseStdout = 0
		}
	}
}

func LoadFromEnv() {
	var e string

	e = os.Getenv("KLOG_SHORT_PATH")
	if e == "1" {
		Conf.ShortPath = 1
	} else if e == "0" {
		Conf.ShortPath = 0
	}

	e = os.Getenv("KLOG_DULL")
	if e == "1" {
		Conf.Dull = 1
	} else if e == "0" {
		Conf.Dull = 0
	}

	e = os.Getenv("KLOG_MUTE")
	if e == "1" {
		Conf.Mute = 1
	} else if e == "0" {
		Conf.Mute = 0
	}

	e = os.Getenv("KLOG_USE_STDOUT")
	if e == "1" {
		Conf.UseStdout = 1
	} else if e == "0" {
		Conf.UseStdout = 0
	}
}

func LoadFromArg() {
	for _, argv := range os.Args {
		if strings.HasPrefix(argv, "--klog-shortPath=") {
			f := argv[17:]
			if f == "1" {
				Conf.ShortPath = 1
			} else if f == "0" {
				Conf.ShortPath = 0
			}
		}
		if strings.HasPrefix(argv, "--klog-dull=") {
			f := argv[12:]
			if f == "1" {
				Conf.Dull = 1
			} else if f == "0" {
				Conf.Dull = 0
			}
		}

		if strings.HasPrefix(argv, "--klog-mute=") {
			f := argv[12:]
			if f == "1" {
				Conf.Mute = 1
			} else if f == "0" {
				Conf.Mute = 0
			}
		}

		if strings.HasPrefix(argv, "--klog-useStdout=") {
			f := argv[17:]
			if f == "1" {
				Conf.UseStdout = 1
			} else if f == "0" {
				Conf.UseStdout = 0
			}
		}
	}
}
