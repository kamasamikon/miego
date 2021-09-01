package pf

import (
	"github.com/kamasamikon/miego/atox"
	// "github.com/kamasamikon/miego/klog"
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
	var js float64
	var jkLo float64
	var jkHi float64
	var zcLo float64
	var zcHi float64
	var ch float64

	/// case 1:
	js = -50
	jkLo = -25
	jkHi = 150
	zcLo = 175
	zcHi = 350
	ch = 375

	slMapQiuJing["1"] = ScoreLineNew(
		js/100, 10, // *近视*
		jkLo/100, 20, // *较快下*
		jkHi/100, 60, // *较快上*
		zcLo/100, 70, // *正常下*
		(zcLo+zcHi)/200, 100, // 正常中 = 100
		zcHi/100, 70, // *正常上*
		ch/100, 50, // *迟缓*
		(ch+50)/100, 30, // 迟缓 + 50
	)

	/// case 2:
	js = -50
	jkLo = -25
	jkHi = 150
	zcLo = 175
	zcHi = 300
	ch = 325
	slMapQiuJing["2"] = ScoreLineNew(
		js/100, 10, // *近视*
		jkLo/100, 20, // *较快下*
		jkHi/100, 60, // *较快上*
		zcLo/100, 70, // *正常下*
		(zcLo+zcHi)/200, 100, // 正常中 = 100
		zcHi/100, 70, // *正常上*
		ch/100, 50, // *迟缓*
		(ch+50)/100, 30, // 迟缓 + 50
	)

	/// case 3:
	js = -50
	jkLo = -25
	jkHi = 150
	zcLo = 175
	zcHi = 300
	ch = 325
	slMapQiuJing["3"] = ScoreLineNew(
		js/100, 10, // *近视*
		jkLo/100, 20, // *较快下*
		jkHi/100, 60, // *较快上*
		zcLo/100, 70, // *正常下*
		(zcLo+zcHi)/200, 100, // 正常中 = 100
		zcHi/100, 70, // *正常上*
		ch/100, 50, // *迟缓*
		(ch+50)/100, 30, // 迟缓 + 50
	)

	/// case 4:
	js = -50
	jkLo = -25
	jkHi = 150
	zcLo = 175
	zcHi = 225
	ch = 250
	slMapQiuJing["4"] = ScoreLineNew(
		js/100, 10, // *近视*
		jkLo/100, 20, // *较快下*
		jkHi/100, 60, // *较快上*
		zcLo/100, 70, // *正常下*
		(zcLo+zcHi)/200, 100, // 正常中 = 100
		zcHi/100, 70, // *正常上*
		ch/100, 50, // *迟缓*
		(ch+50)/100, 30, // 迟缓 + 50
	)

	/// case 5:
	js = -50
	jkLo = -25
	jkHi = 150
	zcLo = 175
	zcHi = 225
	ch = 250
	slMapQiuJing["5"] = ScoreLineNew(
		js/100, 10, // *近视*
		jkLo/100, 20, // *较快下*
		jkHi/100, 60, // *较快上*
		zcLo/100, 70, // *正常下*
		(zcLo+zcHi)/200, 100, // 正常中 = 100
		zcHi/100, 70, // *正常上*
		ch/100, 50, // *迟缓*
		(ch+50)/100, 30, // 迟缓 + 50
	)

	/// case 6:
	js = -50
	jkLo = -25
	jkHi = 125
	zcLo = 150
	zcHi = 225
	ch = 250
	slMapQiuJing["6"] = ScoreLineNew(
		js/100, 10, // *近视*
		jkLo/100, 20, // *较快下*
		jkHi/100, 60, // *较快上*
		zcLo/100, 70, // *正常下*
		(zcLo+zcHi)/200, 100, // 正常中 = 100
		zcHi/100, 70, // *正常上*
		ch/100, 50, // *迟缓*
		(ch+50)/100, 30, // 迟缓 + 50
	)

	/// case 7:
	js = -50
	jkLo = -25
	jkHi = 125
	zcLo = 150
	zcHi = 225
	ch = 250
	slMapQiuJing["7"] = ScoreLineNew(
		js/100, 10, // *近视*
		jkLo/100, 20, // *较快下*
		jkHi/100, 60, // *较快上*
		zcLo/100, 70, // *正常下*
		(zcLo+zcHi)/200, 100, // 正常中 = 100
		zcHi/100, 70, // *正常上*
		ch/100, 50, // *迟缓*
		(ch+50)/100, 30, // 迟缓 + 50
	)

	/// case 8:
	js = -50
	jkLo = -25
	jkHi = 100
	zcLo = 125
	zcHi = 225
	ch = 250
	slMapQiuJing["8"] = ScoreLineNew(
		js/100, 10, // *近视*
		jkLo/100, 20, // *较快下*
		jkHi/100, 60, // *较快上*
		zcLo/100, 70, // *正常下*
		(zcLo+zcHi)/200, 100, // 正常中 = 100
		zcHi/100, 70, // *正常上*
		ch/100, 50, // *迟缓*
		(ch+50)/100, 30, // 迟缓 + 50
	)

	/// case 9:
	js = -50
	jkLo = -25
	jkHi = 75
	zcLo = 100
	zcHi = 200
	ch = 225
	slMapQiuJing["9"] = ScoreLineNew(
		js/100, 10, // *近视*
		jkLo/100, 20, // *较快下*
		jkHi/100, 60, // *较快上*
		zcLo/100, 70, // *正常下*
		(zcLo+zcHi)/200, 100, // 正常中 = 100
		zcHi/100, 70, // *正常上*
		ch/100, 50, // *迟缓*
		(ch+50)/100, 30, // 迟缓 + 50
	)

	/// case 10:
	js = -50
	jkLo = -25
	jkHi = 50
	zcLo = 75
	zcHi = 175
	ch = 200
	slMapQiuJing["10"] = ScoreLineNew(
		js/100, 10, // *近视*
		jkLo/100, 20, // *较快下*
		jkHi/100, 60, // *较快上*
		zcLo/100, 70, // *正常下*
		(zcLo+zcHi)/200, 100, // 正常中 = 100
		zcHi/100, 70, // *正常上*
		ch/100, 50, // *迟缓*
		(ch+50)/100, 30, // 迟缓 + 50
	)

	/// case 11:
	js = -50
	jkLo = -25
	jkHi = 25
	zcLo = 50
	zcHi = 150
	ch = 175
	slMapQiuJing["11"] = ScoreLineNew(
		js/100, 10, // *近视*
		jkLo/100, 20, // *较快下*
		jkHi/100, 60, // *较快上*
		zcLo/100, 70, // *正常下*
		(zcLo+zcHi)/200, 100, // 正常中 = 100
		zcHi/100, 70, // *正常上*
		ch/100, 50, // *迟缓*
		(ch+50)/100, 30, // 迟缓 + 50
	)
}

// s, _ := slMapQiuJing.Score(fInt)
func XXX_QuGuangQiuJing(Gender string, fAge float32, vStr string) int {
	xAge := ""
	switch {
	case fAge >= 11:
		xAge = "11"

	case fAge >= 10:
		xAge = "10"

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

	case fAge >= 3:
		xAge = "3"

	default:
		xAge = "3"
	}

	sl, ok := slMapQiuJing[xAge]
	// klog.Dump(xAge)
	// klog.Dump(sl)
	if ok {
		fInt := atox.Float(vStr, 0)
		s, _ := sl.Score(fInt)
		return int((s+5)/10) * 10
	}

	return 10
}
