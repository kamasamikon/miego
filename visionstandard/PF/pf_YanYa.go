package pf

import (
	"miego/atox"
)

//
// YanYa
//

/*
偏低＜10mmHg
正常10~21
偏高＞21

100分=14-17，增加或减少1减10分
*/
var slMapYanYa *ScoreLine

func init() {
	slMapYanYa = ScoreLineNew(
		// 其他 1 对应 10分
		4, 0,

		// 14 ~ 17 = 100
		14, 100,
		17, 100,

		// 其他 1 对应 10分
		27, 0,
	)
}
func XXX_YanYa(vStr string) int {
	fInt := atox.Float(vStr, 0)
	s, _ := slMapYanYa.Score(fInt)
	return int((s+5)/10) * 10
}
