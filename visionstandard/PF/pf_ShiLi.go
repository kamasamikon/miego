package pf

import (
	"github.com/kamasamikon/miego/atox"
)

//
// ShiLi
//
// < 10 = 偏低
// 10 ~ 21 = 正常
// > 21 = 偏高
//
// 14 ~ 17 = 100
// 其他 1 对应 10分
// V: ScoreLine
var slMapShiLi *ScoreLine

func init() {
	slMapShiLi = ScoreLineNew(
		// 其他 1 对应 10分
		4, 0,

		// 14 ~ 17 = 100
		14, 100,
		17, 100,

		// 其他 1 对应 10分
		27, 0,
	)
}
func XXX_ShiLi(vStr string, fAge float64) int {
	fInt := atox.Float(vStr, 0)
	s, _ := slMapShiLi.Score(fInt)
	return int(s)
}
