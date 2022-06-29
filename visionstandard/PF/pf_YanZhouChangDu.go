package pf

import (
	"fmt"

	"github.com/kamasamikon/miego/atox"
)

//
// YanZhouChangDu
//
// K: "%s_%01f" % (Gender, fAge)
// V: ScoreLine
var slMapYanZhouChangDu map[string]*ScoreLine

func init() {
	slMapYanZhouChangDu = make(map[string]*ScoreLine)

	// (25, ~) = Xxx = 10
	// (22.67, 25) = GaoWei = 10~50
	// (20.5, 22.67) = DiWei = 60~90
	// (~, 20.5) = ZhengChang = 100

	slMapYanZhouChangDu["0_4"] = ScoreLineNew(
		0, 100, // 固定值
		20.5, 100,
		22.67, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["0_5"] = ScoreLineNew(
		0, 100, // 固定值
		21.10, 100,
		23.03, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["0_6"] = ScoreLineNew(
		0, 100, // 固定值
		21.10, 100,
		23.33, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["0_7"] = ScoreLineNew(
		0, 100, // 固定值
		21.50, 100,
		23.50, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["0_8"] = ScoreLineNew(
		0, 100, // 固定值
		21.50, 100,
		23.70, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["0_9"] = ScoreLineNew(
		0, 100, // 固定值
		22.00, 100,
		24.00, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["0_10"] = ScoreLineNew(
		0, 100, // 固定值
		22.00, 100,
		24.00, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["0_11"] = ScoreLineNew(
		0, 100, // 固定值
		22.40, 100,
		24.00, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["0_12"] = ScoreLineNew(
		0, 100, // 固定值
		22.60, 100,
		24.00, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["0_13"] = ScoreLineNew(
		0, 100, // 固定值
		22.90, 100,
		24.10, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["0_14"] = ScoreLineNew(
		0, 100, // 固定值
		23.20, 100,
		24.2, 60,
		25, 10, // 固定值
	)

	slMapYanZhouChangDu["1_4"] = ScoreLineNew(
		0, 100, // 固定值
		20.5, 100,
		22.67, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["1_5"] = ScoreLineNew(
		0, 100, // 固定值
		21.10, 100,
		23.03, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["1_6"] = ScoreLineNew(
		0, 100, // 固定值
		21.10, 100,
		23.33, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["1_7"] = ScoreLineNew(
		0, 100, // 固定值
		21.50, 100,
		23.50, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["1_8"] = ScoreLineNew(
		0, 100, // 固定值
		21.50, 100,
		23.70, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["1_9"] = ScoreLineNew(
		0, 100, // 固定值
		22.00, 100,
		24.00, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["1_10"] = ScoreLineNew(
		0, 100, // 固定值
		22.00, 100,
		24.00, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["1_11"] = ScoreLineNew(
		0, 100, // 固定值
		22.40, 100,
		24.00, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["1_12"] = ScoreLineNew(
		0, 100, // 固定值
		22.60, 100,
		24.00, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["1_13"] = ScoreLineNew(
		0, 100, // 固定值
		22.90, 100,
		24.10, 60,
		25, 10, // 固定值
	)
	slMapYanZhouChangDu["1_14"] = ScoreLineNew(
		0, 100, // 固定值
		23.20, 100,
		24.2, 60,
		25, 10, // 固定值
	)

}

func XXX_YanZhouChangDu(Gender string, fAge float32, vStr string) int {
	fInt := atox.Float(vStr, 0)

	xAge := ""
	switch {
	case fAge >= 14:
		xAge = "14"

	case fAge >= 13:
		xAge = "13"

	case fAge >= 12:
		xAge = "12"

	case fAge >= 11:
		xAge = "11"

	case fAge >= 9:
		xAge = "9"

	case fAge >= 8:
		xAge = "8"

	case fAge >= 7:
		xAge = "7"

	case fAge >= 6:
		xAge = "6"

	case fAge >= 5:
		xAge = "5"

	case fAge >= 4:
		xAge = "4"

	default:
		xAge = "4"
	}

	Key := fmt.Sprintf("%s_%s", Gender, xAge)
	sl, ok := slMapYanZhouChangDu[Key]
	if ok {
		s, _ := sl.Score(fInt)
		return int((s+5)/10) * 10
	}
	return 10
}
