package pf

import (
	"fmt"

	"miego/atox"
	"miego/klog"
)

// BMI
//
// K: "%s_%01f" % (Gender, fAge)
// V: ScoreLine
var slMapBMI map[string]*ScoreLine

func init() {
	slMapBMI = make(map[string]*ScoreLine)

	slMapBMI["0_4.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_4.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_5.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_5.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_6.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_6.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_7.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_7.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_8.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_8.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_9.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_9.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_10.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_10.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_11.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_11.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_12.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_12.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_13.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_13.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_14.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["0_14.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)

	slMapBMI["1_4.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_4.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_5.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_5.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_6.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_6.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_7.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_7.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_8.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_8.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_9.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_9.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_10.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_10.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_11.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_11.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_12.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_12.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_13.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_13.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_14.0"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)
	slMapBMI["1_14.5"] = ScoreLineNew(
		12.0, 20,
		14.71, 50,
		21.71, 90,
		24.40, 50,
		26.00, 10,
	)

	// < 14.71 = 低体重
	// 14.71 ~ 21.71 = 正常
	// 21.71 ~ 24.41 = 超重
	// > 24.41 = 肥胖
	//
	// 正常 = 100
	// 肥胖 = 10
	// 超重、低体重 = 20~60
	slMapBMI["default"] = ScoreLineNew(
		// 超重、低体重 = 20~60
		14.0, 20,

		// 正常 = 100
		14.71, 100,
		21.71, 100,

		// 肥胖 = 10
		24.41, 10,
	)
}

func XXX_BMI(Gender string, fAge float32, Weight string, Height string) int {
	iWeight := atox.Float(Weight, 1)
	iHeight := atox.Float(Height, 1)

	klog.Dump(Gender, "Gender: ")
	klog.Dump(fAge, "fAge: ")
	klog.Dump(iWeight, "iWeight")
	klog.Dump(iHeight, "iHeight")

	// XXX: FIXME: 这里需要计算年龄、性别、体重、身高
	BMI := iWeight * 10000 / iHeight / iHeight
	klog.Dump(BMI, "BMI")

	xAge := ""
	switch {
	case fAge >= 14.5:
		xAge = "14.5"
	case fAge >= 14:
		xAge = "14.0"

	case fAge >= 13.5:
		xAge = "13.5"
	case fAge >= 13:
		xAge = "13.0"

	case fAge >= 12.5:
		xAge = "12.5"
	case fAge >= 12:
		xAge = "12.0"

	case fAge >= 11.5:
		xAge = "11.5"
	case fAge >= 11:
		xAge = "11.0"

	case fAge >= 9.5:
		xAge = "9.5"
	case fAge >= 9:
		xAge = "9.0"

	case fAge >= 8.5:
		xAge = "8.5"
	case fAge >= 8:
		xAge = "8.0"

	case fAge >= 7.5:
		xAge = "7.5"
	case fAge >= 7:
		xAge = "7.0"

	case fAge >= 6.5:
		xAge = "6.5"
	case fAge >= 6:
		xAge = "6.0"

	case fAge >= 5.5:
		xAge = "5.5"
	case fAge >= 5:
		xAge = "5.0"

	case fAge >= 4.5:
		xAge = "4.5"
	case fAge >= 4:
		xAge = "4.0"

	default:
		xAge = "4.0"
	}

	Key := fmt.Sprintf("%s_%s", Gender, xAge)
	Key = "default"
	sl, ok := slMapBMI[Key]
	klog.Dump(sl)
	if ok {
		s, _ := sl.Score(BMI)
		return int((s+5)/10) * 10
	}
	return 10
}
