package roast

import (
	"github.com/kamasamikon/miego/xmap"
	"github.com/kamasamikon/miego/xtime"
	"github.com/kamasamikon/miego/xvx/nvn"
)

func Str__Time(mp xmap.Map, field string) {
	mp.Put(field+"__Str", xtime.NumToStr(field))
}
func Str__Nvn(mp xmap.Map, field string, nvnClass string) {
	mp.Put(field+"__Str", nvn.S(field, nvnClass))
}
