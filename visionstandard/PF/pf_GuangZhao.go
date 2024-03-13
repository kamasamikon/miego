package pf

import (
	"miego/atox"
)

// GuangZhao
//
// < 1 = 10
// 1 ~ 2 = 50
// > 2 = 100
//
// 14 ~ 17 = 100
// 其他 1 对应 10分
// V: ScoreLine
var slMapGuangZhao *ScoreLine

func init() {
	slMapGuangZhao = ScoreLineNew(
		1, 10,
		1.9, 20,
		1.95, 30,
		2, 100,
	)
}
func XXX_GuangZhao(vStr string) int {
	fInt := atox.Float(vStr, 0)
	s, _ := slMapGuangZhao.Score(fInt)
	return int((s+5)/10) * 10
}
