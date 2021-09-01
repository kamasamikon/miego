package pf

import (
	"github.com/kamasamikon/miego/atox"
)

//
// JMQL
//
// 理想 = 100
// 其他 = 每10分 = 0.15,
var slMapJMQL map[string]*ScoreLine

func init() {
	slMapJMQL = make(map[string]*ScoreLine)

	slMapJMQL["4"] = ScoreLineNew(
		42.9, 10,
		43.5, 100,
		44.5, 100,
		45.1, 10,
	)

	slMapJMQL["5"] = ScoreLineNew(
		41.4, 10,
		43, 100,
		44, 100,
		44.6, 10,
	)

	slMapJMQL["6"] = ScoreLineNew(
		41.9, 10,
		42.5, 100,
		43.5, 100,
		44.1, 10,
	)

	slMapJMQL["7"] = ScoreLineNew(
		41.9, 10,
		42.5, 100,
		43.5, 100,
		44.1, 10,
	)

	slMapJMQL["8"] = ScoreLineNew(
		41.9, 10,
		42.5, 100,
		43.5, 100,
		44.1, 10,
	)

	slMapJMQL["9"] = ScoreLineNew(
		41.9, 10,
		42.5, 100,
		43.5, 100,
		44.1, 10,
	)

	slMapJMQL["10"] = ScoreLineNew(
		41.9, 10,
		42.5, 100,
		43.5, 100,
		44.1, 10,
	)

	slMapJMQL["11"] = ScoreLineNew(
		41.9, 10,
		42.5, 100,
		43.5, 100,
		44.1, 10,
	)

}
func XXX_JMQL(vStr string, fAge float32) int {
	fInt := atox.Float(vStr, 0)

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

	default:
		xAge = "4"
	}

	sl, ok := slMapJMQL[xAge]
	if ok {
		s, _ := sl.Score(fInt)
		return int((s+5)/10) * 10
	}
	return 10

}
