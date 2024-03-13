package roast

import (
	"miego/conf"
)

var Conf struct {
	Noisy            bool // call klog
	NotifyQueryError bool // call wxcard
}

func init() {
	conf.OnReady(func() {
		Conf.Noisy = conf.Bool(false, "b:/roast/Noisy")
		Conf.NotifyQueryError = conf.Bool(false, "b:/roast/NotifyQueryError")
	})
}
