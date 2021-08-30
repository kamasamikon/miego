package pf

import (
	"fmt"

	"github.com/kamasamikon/miego/atox"
	"github.com/kamasamikon/miego/klog"
)

func To5X(i string) string {
	switch i {
	case "0.1":
		return "4.0"
	case "0.12":
		return "4.1"
	case "0.15":
		return "4.2"
	case "0.2":
		return "4.3"
	case "0.25":
		return "4.4"
	case "0.3":
		return "4.5"
	case "0.4":
		return "4.6"
	case "0.5":
		return "4.7"
	case "0.6":
		return "4.8"
	case "0.8":
		return "4.9"
	case "1.0":
		return "5.0"
	case "1.2":
		return "5.1"
	case "1.5":
		return "5.2"
	case "2.0":
		return "5.3"
	}

	klog.Dump(i)
	num := atox.Float(i, 0)
	if num == 0 {
		klog.D("xxxx")
		return "5.0"
	}

	if num >= 650 {
		return "4.0"
	}
	if num >= 600 {
		return "4.1"
	}
	if num >= 500 {
		return "4.2"
	}
	if num >= 450 {
		return "4.3"
	}
	if num >= 400 {
		return "4.4"
	}
	if num >= 350 {
		return "4.5"
	}
	if num >= 250 {
		return "4.6"
	}
	if num >= 200 {
		return "4.7"
	}
	if num >= 150 {
		return "4.8"
	}
	if num >= 100 {
		return "4.9"
	}

	return i
}

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

	// XXX
	if score < 10 {
		score = 10
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
	//
	// 视力转成屈光度
	//
	// 1.0 对数
	// 4.9 指数
	// 200 度数
	//

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

	// XXX
	if score < 10 {
		score = 10
	}

	return int(score)
}

// 14.71 21.71 24.41

// BMI评分计算
// if (bmi >=22) score = -10 * bmi + 320
// if (bmi < 22) score = 10 * bmi - 120
// if (bmi < 18) score = 10

func BMI(Gender string, fAge float32, Weight string, Height string) int {
	// Age := int(fAge)

	iWeight := atox.Float(Weight, 1)
	iHeight := atox.Float(Height, 1)
	klog.Dump(iWeight)
	klog.Dump(iHeight)

	// XXX: FIXME: 这里需要计算年龄、性别、体重、身高
	bmi := iWeight * 10000 / iHeight / iHeight
	klog.Dump(bmi)

	score := 0

	if bmi >= 22 {
		score = int(-10*bmi + 320)
	}
	if bmi < 22 {
		score = int(10*bmi - 120)
	}
	// if bmi < 18 {
	// score = 10
	// }

	if score < 0 {
		score = 0
	}
	if score > 100 {
		score = 100
	}

	// XXX
	if score < 10 {
		score = 10
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

	// XXX
	if score < 10 {
		score = 10
	}

	return int(score)
}

func ShiLi(vStr string, fAge float32) int {
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
	klog.Dump(vInt, "vInt: ")
	klog.Dump(Age, "Age: ")

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
	klog.Dump(score)
	score /= 1.5
	klog.Dump(score)

	if score > 95 {
		score = 95
	}
	if score < 10 {
		score = 10
	}

	return int(score)
}

// 角膜曲率计
// 常数：R41=43.5 R42=44.5 R51=43 R52=44 R61=42.5 R62=43.5 ...
// 对于x岁孩子，角膜曲率测得为RR，通过与Rx进行对比计算：
// 5岁孩子，曲率为43.8，得分score计算
// if (RR >= R51 && RR <= R52) score = 90
// if (RR > R52) score = 95
// if (RR < R51) score = 90 - (R51 - RR) * 100
func JMQL(vStr string, fAge float32) int {
	Age := int(fAge)

	if Age < 4 {
		Age = 4
	}
	if Age > 11 {
		Age = 11
	}

	vInt := atox.Float(vStr, 0)

	var x1 float64
	var x2 float64

	switch Age {
	case 4:
		x1 = 43.50
		x2 = 44.50

	case 5:
		x1 = 43.00
		x2 = 44.00

	case 6:
		x1 = 42.50
		x2 = 43.50

	case 7:
		x1 = 42.50
		x2 = 43.50

	case 8:
		x1 = 42.50
		x2 = 43.50

	case 9:
		x1 = 42.50
		x2 = 43.50

	case 10:
		x1 = 42.50
		x2 = 43.50

	case 11:
		x1 = 42.50
		x2 = 43.50
	}

	klog.Dump(fAge, "fAge")
	klog.Dump(Age, "Age")
	klog.Dump(x1, "x1")
	klog.Dump(x2, "x2")
	klog.Dump(vInt, "vInt")

	// 斜率（incline）= 分数/长度
	// 分数: 10(上限) - 50(比较差)
	// 中位数 = 95
	// 边界 = 50
	middle := (x2 + x1) / 2
	incline := 22.5 / ((x2 - x1) / 2)
	klog.Dump(incline, "incline")
	klog.Dump(middle, "middle")

	var score float64
	if vInt > middle {
		score = 95 - (vInt-middle)*incline
	} else {
		score = 95 + (vInt-middle)*incline
	}

	klog.Dump(score)

	if score > 95 {
		score = 95
	}
	if score < 10 {
		score = 10
	}

	return int(score)
}

// 常数：Q41=2 Q42=3 Q51=1.5 Q52=2.5 Q61=1 Q62=2 Q71=0.5 Q72=1.5 Q81=0 Q82=1.0
// Q91=0 Q92=0.5 Qa1=0 Qa2=0.5 Qb1=0 Qb2=0.5
// 对于x岁孩子，屈光度球镜测得为Qr，通过与Qx进行对比计算：
// 5岁孩子，屈光度球镜为2，得分score计算
// if (QR >= Q51 && QR <= Q52) score = 90;
// if (QR > Q52) score = 95;
// if (QR < Q51) score = 90 - (Q51 - QR) * 100;
func QuGuangQiuJing(Gender string, fAge float32, vStr string) int {
	Age := int(fAge)

	if Age < 4 {
		Age = 4
	}
	if Age > 11 {
		Age = 11
	}

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
		score = 90
		klog.Dump(score)
	}

	if vInt > x2 {
		score = 95
		klog.Dump(score)
	}

	if vInt < x1 {
		score = 90 - (x1-vInt)/100*100
		klog.Dump(score)
	}

	if score < 0 {
		score = 0
		klog.Dump(score)
	}
	if score > 100 {
		score = 100
		klog.Dump(score)
	}

	// XXX
	if score < 10 {
		score = 10
		klog.Dump(score)
	}

	klog.Dump(score)
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

func YanZhouChangDu(Gender string, fAge float32, vStr string) int {
	Age := int(fAge)

	vInt := float32(atox.Float(vStr, 0))

	// var x0 int
	var x1 float32
	var x2 float32

	if Age > 11 {
		Age = 11
	}
	if Age < 4 {
		Age = 4
	}

	if Gender == "0" {
		// 女
		switch Age {
		case 4:
			x1 = 22.06
			x2 = 22.13

		case 5:
			x1 = 23.40
			x2 = 23.46

		case 6:
			x1 = 23.71
			x2 = 23.78

		case 7:
			x1 = 22.92
			x2 = 23.09

		case 8:
			x1 = 23.18
			x2 = 23.34

		case 9:
			x1 = 23.52
			x2 = 23.61

		case 10:
			x1 = 23.52
			x2 = 23.87

		case 11:
			x1 = 23.52
			x2 = 24.07

		}

	} else {
		// 男
		switch Age {
		case 4:
			x1 = 22.67
			x2 = 22.70

		case 5:
			x1 = 23.03
			x2 = 23.05

		case 6:
			x1 = 23.33
			x2 = 23.37

		case 7:
			x1 = 23.50
			x2 = 23.67

		case 8:
			x1 = 23.70
			x2 = 23.90

		case 9:
			x1 = 24.01
			x2 = 24.13

		case 10:
			x1 = 24.01
			x2 = 24.40

		case 11:
			x1 = 24.01
			x2 = 24.41
		}
	}

	// > 上限 = 高危
	// > 下线 = 中卫
	// < 下线 = 正常

	// 越长越差
	// x2 最差
	// x1 比较差

	// klog.Dump(x1)
	// klog.Dump(x2)
	// klog.Dump(vInt)
	// 斜率（incline）= 分数/长度
	// 分数: 10(上限) - 50(比较差)
	incline := 40 / (x2 - x1)
	// klog.Dump(incline)

	// 以50分线为准
	score := 50 - (vInt-x1)*incline
	// klog.Dump(score)

	if score > 95 {
		score = 95
	}
	if score < 10 {
		score = 10
	}

	return int(score)

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

//
// 评分折线，由多个分数的点串起来的折线。（最好是曲线，但不知道怎么实现好）
//
type ScoreItem struct {
	Start float64
	Score float64
}

type ScoreLine struct {
	items []*ScoreItem
}

// v: 数值 VS 对应的分数
// 要求按照分数从小到大排序
func ScoreLineNew(v ...float64) *ScoreLine {
	Line := &ScoreLine{}

	for i := 0; i < len(v)/2; i++ {
		Start := v[2*i]
		Score := v[2*i+1]

		Line.items = append(Line.items, &ScoreItem{Start, Score})
	}

	return Line
}

// 输入一个数，返回这数对应的评分
func (sl *ScoreLine) Score(value float64) (score float64, kind int) {
	items := sl.items

	IndexNext := -1
	for i := range items {
		if value < items[i].Start {
			IndexNext = i
			break
		}
	}

	switch IndexNext {
	case 0:
		kind = -1

		score = items[0].Score

	case -1:
		kind = 1

		Size := len(items)
		score = items[Size-1].Score

	default:
		kind = 0

		ItemNext := items[IndexNext]
		ItemCurr := items[IndexNext-1]
		incline := (ItemNext.Score - ItemCurr.Score) / (ItemNext.Start - ItemCurr.Start)
		score = ItemCurr.Score + incline*(value-ItemCurr.Start)
	}

	return score, kind
}
