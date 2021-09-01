package pf

import (
	"github.com/kamasamikon/miego/atox"
)

//
// ShiLi
//

/*
正常60-100，低常10-50，10分≤4.0，100分≥5.2
*/

// 4.0 = 10
// 4.9 = 10 // 两头好算，中间如何取点？
// 5.2 = 100

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
func XXX_ShiLi(vStr string, fAge float32) int {
	Age := int(fAge)

	vStr = To5X(vStr)

	if Age < 4 {
		Age = 4
	}
	if Age > 11 {
		Age = 11
	}

	S4 := 4.8
	S5 := 4.85
	S6 := 4.9
	S7 := 4.95
	S8 := 5.0
	S9 := 5.0
	S10 := 5.0
	S11 := 5.0

	vInt := atox.Float(vStr, 0)

	// 100 -> 100
	// 90 -> 60

	//

	var x float64

	switch Age {
	case 4:
		x = S4

	case 5:
		x = S5

	case 6:
		x = S6

	case 7:
		x = S7

	case 8:
		x = S8

	case 9:
		x = S9

	case 10:
		x = S10

	case 11:
		x = S11
	}

	score := 10 + (vInt-4.0)*(80/(x-4.0))
	score /= 1.5

	if score > 95 {
		score = 95
	}
	if score < 10 {
		score = 10
	}

	s := score
	return int((s+5)/10) * 10
}
