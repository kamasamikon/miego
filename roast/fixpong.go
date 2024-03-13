package roast

import (
	"miego/xmap"
)

func FixPongs(pongs []xmap.Map, offset uint, allCount int) {
	for i := range pongs {
		pongs[i]["Index"] = i
		pongs[i]["AllIndex"] = int(offset) + i
		pongs[i]["AllCount"] = allCount
	}
}
