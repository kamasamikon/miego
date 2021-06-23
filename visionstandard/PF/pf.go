package pf

import (
	"fmt"

	"github.com/kamasamikon/miego/atox"
	"github.com/kamasamikon/miego/klog"
)

// 眼健康综合
// 综合视力风险
// 近视防控等级
// 体重
// 遗传
// 光照
// 视力
// 屈光

// 眼压p，单位：mmHg
// if (10 <= p <= 20) score = 100;
// if (p < 10) score = 50 * p - 400;
// if (p > 20) score = -50 * p + 1100;
// if (score < 0) score = 0;
func YanYa(vStr string) int {
	vInt := int(atox.Float(vStr, 0))

	var score int
	if 10 <= vInt && vInt <= 20 {
		score = 100
	}
	if vInt < 10 {
		score = 50*vInt - 400
	}
	if vInt > 20 {
		score = -50*vInt + 1100
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// 遗传评分计算
// d1是父亲的屈光度，d2是母亲的屈光度
// score = -7.5 * (d1 + d2) + 100；
// if (score < 0) score = 0
// if (score >100) score = 100
// F/M: 度数
func YiChuan(F string, M string) int {
	if F == "" {
		F = "0"
	}
	if M == "" {
		M = "0"
	}

	iF := atox.Int(F, 0)
	iM := atox.Int(M, 0)
	score := 100 - (75 * (iF + iM) / 100 / 10)

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return int(score)
}

// BMI评分计算
// if (bmi >=22) score = -10 * bmi + 320
// if (bmi < 22) score = 10 * bmi - 120
// if (bmi < 18) score = 10

func BMI(Gender string, Age int, Weight string, Height string) int {
	iWeight := atox.Int(Weight, 1)
	iHeight := atox.Int(Height, 1)
	klog.Dump(iWeight)
	klog.Dump(iHeight)
	// XXX: FIXME: 这里需要计算年龄、性别、体重、身高
	bmi := iWeight * 10000 / iHeight / iHeight
	klog.Dump(bmi)

	score := 0

	if bmi >= 22 {
		score = -10*bmi + 320
	}
	if bmi < 22 {
		score = 10*bmi - 120
	}
	if bmi < 18 {
		score = 10
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// 光照评分计算
// 全天光照时间（含课间休息室外时间）：sun 单位：分钟
func GuangZhao(hours string) int {
	sun := atox.Float(hours, 0)
	sun *= 60

	score := 0.5*sun + 30

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return int(score)
}

func ShiLi(vStr string, Age int) int {
	S4 := 60
	S5 := 70
	S6 := 80
	S7 := 90
	S8 := 100
	S9 := 100
	S10 := 100
	S11 := 100

	vInt := int(atox.Float(vStr, 0) * 100)

	var x int

	score := 0
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

	if vInt >= x {
		score = 90
	}

	if vInt < x {
		score = 90 - (vInt-x)*200
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// 角膜曲率计
// 常数：R41=43.5 R42=44.5 R51=43 R52=44 R61=42.5 R62=43.5 ...
// 对于x岁孩子，角膜曲率测得为RR，通过与Rx进行对比计算：
// 5岁孩子，曲率为43.8，得分score计算
// if (RR >= R51 && RR <= R52) score = 90
// if (RR > R52) score = 95
// if (RR < R51) score = 90 - (R51 - RR) * 100
func JMQL(vStr string, Age int) int {
	vInt := int(atox.Float(vStr, 0) * 100)

	var x1 int
	var x2 int

	score := 0
	switch Age {
	case 4:
		x1 = 4350
		x2 = 4450

	case 5:
		x1 = 4300
		x2 = 4400

	case 6:
		x1 = 4250
		x2 = 4350

	case 7:
		x1 = 4250
		x2 = 4350

	case 8:
		x1 = 4250
		x2 = 4350

	case 9:
		x1 = 4250
		x2 = 4350

	case 10:
		x1 = 4250
		x2 = 4350

	case 11:
		x1 = 4250
		x2 = 4350
	}

	klog.D("vInt:%d, x1:%d, x2:%d", vInt, x1, x2)
	if vInt >= x1 && vInt <= x2 {
		score = 90
	}
	if vInt > x2 {
		score = 95
	}
	if vInt < x1 {
		score = 90 - (x1-vInt)/100*100
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	return score
}

// 常数：Q41=2 Q42=3 Q51=1.5 Q52=2.5 Q61=1 Q62=2 Q71=0.5 Q72=1.5 Q81=0 Q82=1.0
// Q91=0 Q92=0.5 Qa1=0 Qa2=0.5 Qb1=0 Qb2=0.5
// 对于x岁孩子，屈光度球镜测得为Qr，通过与Qx进行对比计算：
// 5岁孩子，屈光度球镜为2，得分score计算
// if (QR >= Q51 && QR <= Q52) score = 90;
// if (QR > Q52) score = 95;
// if (QR < Q51) score = 90 - (Q51 - QR) * 100;
func QuGuangQiuJing(Gender string, Age int, vStr string) int {
	vInt := int(atox.Float(vStr, 0) * 100)

	var x1 int
	var x2 int

	score := 0
	switch Age {
	case 4:
		x1 = 200
		x2 = 300

	case 5:
		x1 = 150
		x2 = 250

	case 6:
		x1 = 100
		x2 = 200

	case 7:
		x1 = 50
		x2 = 150

	case 8:
		x1 = 0
		x2 = 100

	case 9:
		x1 = 0
		x2 = 50

	case 10:
		x1 = 0
		x2 = 50

	case 11:
		x1 = 0
		x2 = 50
	}

	klog.D("vInt:%d, x1:%d, x2:%d", vInt, x1, x2)
	if vInt >= x1 && vInt <= x2 {
		klog.D("")
		score = 90
	}

	if vInt > x2 {
		klog.D("")
		score = 95
	}

	if vInt < x1 {
		klog.D("")
		score = 90 - (x1-vInt)/100*100
	}

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	klog.D("%d", score)
	return score
}

/*
指导建议

同学你的视力健康 [A] ，近视风险 [B] ，考虑到 [C] 明显的父母遗传因素，[D]，[E] ，同时 [F] ，[G] 。

A：

if (视力健康综合评分 >= 70) A = 比较好

else A = 不理想

B：

if (近视风险指数 >= 3) B = 比较大

else B = 不大

C：

if (遗传 >= 60) C = 没有

else B = 有

D：

if (近视风险指数 >= 3) D = 要非常重视近视的防控与治疗

else D = 要保持良好的用眼习惯，防止近视发生

E：

if （光照评分 >= 80） E：请坚持充足的白天户外活动

else E = 请加强白天户外活动，保证每天至少两个小时的白天户外活动

F：

if （BMI >= 70） F：坚持锻炼，保持体型健康

else F = 加强锻炼，增进身体健康

G：

G = 请同学定期做视力健康的检查！
*/

func JianYi(ShiLiZongHe int, JinShiFengXian int, YiChuan int, GuangZhao int, BMI int) string {

	s := `同学你的视力健康 %s ，近视风险 %s ，考虑到 %s 明显的父母遗传因素，%s ，%s ，同时 %s ，%s`

	var A, B, C, D, E, F, G string

	if ShiLiZongHe >= 70 {
		A = "比较好"
	} else {
		A = "不理想"
	}

	if JinShiFengXian >= 3 {
		B = "比较大"
	} else {
		B = "不大"
	}

	if YiChuan >= 60 {
		C = "没有"
	} else {
		C = "有"
	}

	if JinShiFengXian >= 3 {
		D = "要非常重视近视的防控与治疗"
	} else {
		D = "要保持良好的用眼习惯，防止近视发生"
	}

	if GuangZhao >= 80 {
		E = "请坚持充足的白天户外活动"
	} else {
		E = "请加强白天户外活动，保证每天至少两个小时的白天户外活动"
	}

	if BMI >= 70 {
		F = "坚持锻炼，保持体型健康"
	} else {
		F = "加强锻炼，增进身体健康"
	}

	G = "请同学定期做视力健康的检查！"

	return fmt.Sprintf(
		s,
		A, B, C, D, E, F, G,
	)
}

func YanZhouChangDu(Gender string, Age int, vStr string) int {
	vInt := int(atox.Float(vStr, 0))

	// var x0 int
	var x1 int
	var x2 int

	if Gender == "0" {
		// 女
		switch Age {
		case 4:
			x1 = 2206
			x2 = 2213

		case 5:
			x1 = 2340
			x2 = 2346

		case 6:
			x1 = 2371
			x2 = 2378

		case 7:
			x1 = 2292
			x2 = 2309

		case 8:
			x1 = 2318
			x2 = 2334

		case 9:
			x1 = 2352
			x2 = 2361

		case 10:
			x1 = 2352
			x2 = 2387

		case 11:
			x1 = 2352
			x2 = 2407

		}

	} else {
		// 男

		switch Age {
		case 4:
			x1 = 2267
			x2 = 2270

		case 5:

			x1 = 2303
			x2 = 2305

		case 6:

			x1 = 2333
			x2 = 2337

		case 7:

			x1 = 2350
			x2 = 2367

		case 8:

			x1 = 2370
			x2 = 2390

		case 9:

			x1 = 2401
			x2 = 2413

		case 10:

			x1 = 2401
			x2 = 2440

		case 11:

			x1 = 2401
			x2 = 2441
		}

	}

	if vInt < x1 {
		return 20
	}
	if vInt > x2 {
		return 90
	}
	return (vInt-x1)/(x2-x1)*7/10 + 20

	/*

		var s int // score
		x0 = 2*x1 - x2

		x := 20 // 超过某岁的天数

		yc0 := y40 + (y50-y40)*x/365 // 正常
		yc1 := y41 + (y51-y41)*x/365 // 低危
		yc2 := y42 + (y52-y42)*x/365 // 高危

		MS := 90
		MG := 50
		MB := 20

		if yr < yc0 {
			s = MS + (MS-MG)*(yc0-yr)/(yc1-yc0)
		}
		if s > 100 {
			s = 100
		}

		if yr == yc0 {
			s = MS
		}

		if yr > yc0 && yr < yc1 {
			s = MS - (MS-MG)*(yr-yc0)/(yc1-yc0)
		}

		if yr == yc1 {
			s = MG
		}

		if yr > yc1 && yr < yc2 {
			s = MG - (MG-MB)*(yr-yc1)/(yc2-yc1)
		}

		if yr == yc2 {
			s = MB
		}

		if yr > yc2 {
			s = MB - (MG-MB)*(yr-yc2)/(yc2-yc1)
		}
		if s < 0 {
			s = 0
		}
	*/
}
