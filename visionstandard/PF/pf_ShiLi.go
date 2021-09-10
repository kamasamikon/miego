package pf

import (
	"github.com/kamasamikon/miego/atox"
)

//
// ShiLi
//

/*
正常60-100，低常10-50，10分≤4.0，100分≥5.2

5.0 = 80
5.1 = 90
5.2 = 100

4.0 = 10
*/

// 4.0 = 10
// 4.9 = 10 // 两头好算，中间如何取点？
// 5.2 = 100

var slMapShiLi *ScoreLine

func init() {
	slMapShiLi = ScoreLineNew(
		4, 10,

		4.8, 50,

		5.0, 80,
		5.1, 90,

		5.2, 100,
	)
}

func XXX_ShiLi(vStr string, fAge float32) int {
	vStr = To5X(vStr)
	vInt := atox.Float(vStr, 0)

	s, _ := slMapShiLi.Score(vInt)
	return int((s+5)/10) * 10
}
