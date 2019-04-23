package xconf

import (
	"github.com/kamasamikon/miego/conf"
)

func init() {
	conf.Load("./msa.cfg")
	conf.Load("./main.cfg")
	conf.Load("/tmp/conf/main.cfg")
}
