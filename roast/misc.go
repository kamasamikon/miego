package roast

import (
	"github.com/kamasamikon/miego/xmap"
	"github.com/kamasamikon/miego/xtime"
	"github.com/kamasamikon/miego/xvx/nvn"
)

func Str__Time(mp xmap.Map, field string) {
	mp.Put(field+"__Str", xtime.NumToStr(mp.S(field)))
}
func Str__Nvn(mp xmap.Map, field string, nvnClass string) {
	mp.Put(field+"__Str", nvn.S(mp.S(field), nvnClass))
}

/////////////////////////////////////////////////////////////////////////
// For Query:
// 1. PageSize: default to 10
// 2. ID: clear RemAt
// 3. UUID: default to last
func Preset(mp xmap.Map, Prefix string) string {
	var preset string

	// Default to 50 lines
	// mp.SafePut("PageSize", "50")

	// If set ID,
	if mp.Has("ID") {
		preset = ""
	} else if !mp.Has("RemAt") {
		preset = Prefix + ".RemAt = 0"
	}
	if mp.Has("UUID") {
		mp.Put(
			"PageSize", "1",
			"OrderBy", "ID",
			"OrderDir", "desc",
		)
	}
	return preset
}
