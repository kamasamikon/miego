package xconf

import (
	"github.com/kamasamikon/miego/conf"
	"github.com/kamasamikon/miego/klog"
)

func init() {
	files := []string{
		"./msa.cfg",
		"./main.cfg",
		"/tmp/conf/main.cfg",
	}
	for _, f := range files {
		if err := conf.Load(f); err != nil {
			klog.E(err.Error())
		}
	}
}
