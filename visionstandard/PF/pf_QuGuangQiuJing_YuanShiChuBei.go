package pf

import (
	"github.com/kamasamikon/miego/atox"
	"github.com/kamasamikon/miego/klog"
)

//
// QiuJing
//
/*
按10分为1单位，疑似近视10分，屈光发育较快20—60，正常70—100分，屈光发育迟缓50分。
*/

var slMapQiuJing map[string]*ScoreLine

func init() {
	slMapQiuJing = make(map[string]*ScoreLine)

	slMapQiuJing["4"] = ScoreLineNew(
		-0.50, 10, // ** 近视
		1.50, 20, // ** 较快

		1.75, 70, // ** 正常
		2.00, 100, // --- 差值: 中间靠左
		2.50, 50, // ** 迟缓

		2.75, 40, // -- 延续 + 25度
	)

	// 4 和 5 是一样的
	slMapQiuJing["5"] = ScoreLineNew(
		-0.50, 10, // ** 近视
		1.50, 20, // ** 较快

		1.75, 70, // ** 正常
		2.00, 100, // --- 差值: 中间靠左
		2.50, 50, // ** 迟缓

		2.75, 40, // -- 延续 + 25度
	)

	slMapQiuJing["6"] = ScoreLineNew(
		-0.50, 10, // ** 近视
		1.25, 20, // ** 较快

		1.50, 70, // ** 正常
		2.00, 100, // --- 差值: 中间靠左
		2.50, 50, // ** 迟缓

		2.75, 40, // -- 延续 + 25度
	)

	// Same as 6
	slMapQiuJing["7"] = ScoreLineNew(
		-0.50, 10, // ** 近视
		1.25, 20, // ** 较快

		1.50, 70, // ** 正常
		2.00, 100, // --- 差值: 中间靠左
		2.50, 50, // ** 迟缓

		2.75, 40, // -- 延续 + 25度
	)

	slMapQiuJing["8"] = ScoreLineNew(
		-0.50, 10, // ** 近视
		1.00, 20, // ** 较快

		1.25, 70, // ** 正常
		2.00, 100, // --- 差值: 中间靠左
		2.50, 50, // ** 迟缓

		2.75, 40, // -- 延续 + 25度
	)

	slMapQiuJing["9"] = ScoreLineNew(
		-0.50, 10, // ** 近视
		0.75, 20, // ** 较快

		1.00, 70, // ** 正常
		2.00, 100, // --- 差值: 中间靠左
		2.25, 50, // ** 迟缓

		2.50, 40, // -- 延续 + 25度
	)

	// Same as 9
	slMapQiuJing["10"] = ScoreLineNew(
		-0.50, 10, // ** 近视
		0.75, 20, // ** 较快

		1.00, 70, // ** 正常
		2.00, 100, // --- 差值: 中间靠左
		2.25, 50, // ** 迟缓

		2.50, 40, // -- 延续 + 25度
	)

	// Same as 9
	slMapQiuJing["11"] = ScoreLineNew(
		-0.50, 10, // ** 近视
		0.75, 20, // ** 较快

		1.00, 70, // ** 正常
		2.00, 100, // --- 差值: 中间靠左
		2.25, 50, // ** 迟缓

		2.50, 40, // -- 延续 + 25度
	)

}

// s, _ := slMapQiuJing.Score(fInt)
func XXX_QuGuangQiuJing(Gender string, fAge float32, vStr string) int {
	xAge := ""
	switch {
	case fAge >= 11:
		xAge = "11"

	case fAge >= 9.5:
		xAge = "9.5"
	case fAge >= 9:
		xAge = "9"

	case fAge >= 8.5:
		xAge = "8.5"
	case fAge >= 8:
		xAge = "8"

	case fAge >= 7.5:
		xAge = "7.5"
	case fAge >= 7:
		xAge = "7"

	case fAge >= 6.5:
		xAge = "6.5"
	case fAge >= 6:
		xAge = "6"

	case fAge >= 5.5:
		xAge = "5.5"
	case fAge >= 5:
		xAge = "5"

	case fAge >= 4.5:
		xAge = "4.5"
	case fAge >= 4:
		xAge = "4"

	default:
		xAge = "4"
	}

	sl, ok := slMapQiuJing[xAge]
	klog.Dump(xAge)
	klog.Dump(sl)
	if ok {
		fInt := atox.Float(vStr, 0)
		s, _ := sl.Score(fInt)
		return int((s+5)/10) * 10
	}

	return 10
}
