package pf

import (
	"miego/atox"
)

// YiChuan
//
// < 10 = 偏低
// 10 ~ 21 = 正常
// > 21 = 偏高
//
// 14 ~ 17 = 100
// 其他 1 对应 10分
// V: ScoreLine
var slMapYiChuan *ScoreLine

func init() {
	slMapYiChuan = ScoreLineNew(
		// 其他 1 对应 10分
		4, 0,

		// 14 ~ 17 = 100
		14, 100,
		17, 100,

		// 其他 1 对应 10分
		27, 0,
	)
}
func XXX_YiChuan(F string, M string) int {
	fInt := atox.Int(F, 0)
	mInt := atox.Int(M, 0)

	Risk := 0.0
	switch {

	case fInt >= 600:
		switch {
		case mInt >= 600:
			Risk = 50
		case mInt >= 300:
			Risk = 30
		case mInt > 0:
			Risk = 30
		case mInt == 0:
			Risk = 30
		}
	case fInt >= 300:
		switch {
		case mInt >= 600:
			Risk = 30
		case mInt >= 300:
			Risk = 30
		case mInt > 0:
			Risk = 20
		case mInt == 0:
			Risk = 20
		}
	case fInt > 0:
		switch {
		case mInt >= 600:
			Risk = 30
		case mInt >= 300:
			Risk = 20
		case mInt > 0:
			Risk = 20
		case mInt == 0:
			Risk = 10
		}
	case fInt == 0:
		switch {
		case mInt >= 600:
			Risk = 30
		case mInt >= 300:
			Risk = 20
		case mInt > 0:
			Risk = 10
		case mInt == 0:
			Risk = 10
		}
	}

	score := 60
	switch Risk {
	case 0:
		score = 100
	case 10:
		score = 100
	case 20:
		score = 60
	case 30:
		score = 30
	case 50:
		score = 10
	}

	s := score
	return int((s+5)/10) * 10
}
