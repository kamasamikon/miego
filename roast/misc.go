package roast

import (
	"github.com/kamasamikon/miego/xmap"
	"github.com/kamasamikon/miego/xtime"
)

func NewAt__Str(p xmap.Map) {
	p.Put("NewAt__Str", xtime.NumToStr(p.S("NewAt")))
}

func CrtAt__Str(p xmap.Map) {
	p.Put("CrtAt__Str", xtime.NumToStr(p.S("CrtAt")))
}
