package pf

import ()

// 眼健康综合
// 综合视力风险
// 近视防控等级
// 体重
// 遗传
// 光照
// 视力
// 屈光

// 遗传评分计算
// d1是父亲的屈光度，d2是母亲的屈光度
// score = -7.5 * (d1 + d2) + 100；
// if (score < 0) score = 0
// if (score >100) score = 100
func YiChuan(F string, M string) int {
	if F == "" {
		F = "0"
	}
	if M == "" {
		M = "0"
	}

	iF := atox.Int(F, 0)
	iM := atox.Int(M, 0)
	score := -7.5*(iF+iM) + 100

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

// func BMI(Gender string, Age int, Weight int, Height int) int {
func BMI(kg int) int {
	// XXX: FIXME: 这里需要计算年龄、性别、体重、身高
	if bmi >= 22 {
		score = -10*bmi + 320
	}
	if bmi < 22 {
		score = 10*bmi - 120
	}
	if bmi < 18 {
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

	if score > 100 {
		score = 100
	}
	return int(score)
}

func SHILI(fSR string, Age int) int {

	S4 := 60
	S5 := 70
	S6 := 80
	S7 := 90
	S8 := 100
	S9 := 100
	S10 := 100
	S11 := 100

	SR := int(atox.Float(fSR, 0) * 100)
	klog.Dump(SR)

	var Sx int

	score := 0
	switch Age {
	case 4:
		Sx = S4

	case 5:
		Sx = S5

	case 6:
		Sx = S6

	case 7:
		Sx = S7

	case 8:
		Sx = S8

	case 9:
		Sx = S9

	case 10:
		Sx = S10

	case 11:
		Sx = S11
	}

	if SR >= Sx {
		score = 90
	}

	if SR < Sx {
		score = 90 - (SR-Sx)*200
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
func JMQL(fRR string, Age int) int {
	RR := int(atox.Float(fRR, 0) * 100)
	klog.Dump(RR)

	var Rx1 int
	var Rx2 int

	score := 0
	switch Age {
	case 4:
		Rx1 = S4

	case 5:
		Rx1 = S5

	case 6:
		Rx1 = S6

	case 7:
		Rx1 = S7

	case 8:
		Rx1 = S8

	case 9:
		Rx1 = S9

	case 10:
		Rx1 = S10

	case 11:
		Rx1 = S11
	}

	if RR >= Rx1 && RR <= Rx2 {
		score = 90
	}
	if RR > Rx2 {
		score = 95
	}
	if RR < Rx1 {
		score = 90 - (Rx1-RR)*100
	}

	return score
}

// 屈光度球镜评分计算
// 常数：Q41:=2 Q42:=3 Q51:=1.5 Q52:=2.5 Q61:=1 Q62:=2 Q71:=0.5 Q72:=1.5 Q81:=0 Q82:=1.0
func QuGuangQiuJing(Gender string, Age int, Value string) (Result string, Hint string) {
	SR := int(atox.Float(fSR, 0) * 100)
	klog.Dump(SR)

	var Rx1 int
	var Rx2 int

	score := 0
	switch Age {
	case 4:
		Rx1 = S4

	case 5:
		Rx1 = S5

	case 6:
		Rx1 = S6

	case 7:
		Rx1 = S7

	case 8:
		Rx1 = S8

	case 9:
		Rx1 = S9

	case 10:
		Rx1 = S10

	case 11:
		Rx1 = S11
	}

	if QR >= Qx1 && QR <= Qx2 {
		score = 90
	}

	if QR > Qx2 {
		score = 95
	}

	if QR < Qx1 {
		score = 90 - (Qx1-QR)*100
	}
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
		B = "有"
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
