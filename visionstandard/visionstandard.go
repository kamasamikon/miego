package visionstandard

//裸眼和矫正
func Vision(age int) string {
	if age >= 0 && age < 1 {
		return "-"
	} else if age >= 1 && age < 2 {
		return "≥4.2"
	} else if age >= 2 && age < 3 {
		return "≥4.7"
	} else if age >= 3 && age < 4 {
		return "≥4.8"
	} else if age >= 4 && age < 5 {
		return "≥4.9"
	} else {
		return "≥5.0"
	}
}

//球镜
func DS(age int) string {
	if age >= 0 && age < 2 {
		return "+1.75 ~ +3.50"
	} else if age >= 2 && age < 4 {
		return "+1.75 ~ +3.00"
	} else if age >= 4 && age < 6 {
		return "+1.75 ~ +2.25"
	} else if age >= 6 && age < 8 {
		return "+1.50 ~ +2.25"
	} else if age >= 8 && age < 9 {
		return "+1.25 ~ +2.25"
	} else if age >= 9 && age < 10 {
		return "+1.00 ~ +2.00"
	} else if age >= 10 && age < 11 {
		return "+0.75 ~ +1.75"
	} else if age >= 11 && age < 12 {
		return "+0.5 ~ +1.50"
	} else {
		return "0.00 ~ +1.25"
	}
}

//柱镜
func DC(age int) string {
	if age >= 0 && age < 2 {
		return ">-1.50"
	} else {
		return ">-1.00"
	}
}

//水平固视
func GazeH() string {
	return "<8"
}

//垂直固视
func GazeV() string {
	return "<8"
}

//屈光参差(球镜S)
func DSDiff() string {
	return "<1.50"
}

//屈光参差(柱镜C)
func DCDiff() string {
	return "<1.00"
}

//瞳孔直径差
func PSDiff() string {
	return "<1mm"
}

// 近视分级
// 无，轻度、中度、高度
func JinShiChengDu(degree float64) (int, string) {
	// 轻度近视是300度以下近视、
	// 中度近视300-600度之间、
	// 高度近视600-900度之间，
	// 而超高度近视是900度以上的近视。
	switch {
	case degree <= 0:
		return 0, "无"
	case degree < 300:
		return 1, "轻度"
	case degree < 600:
		return 2, "中度"
	case degree < 900:
		return 3, "高度"
	default:
		return 4, "超高度"
	}
}

///////////////////////////////
// 返回具体的数字
//
// 远视储备：屈光：球镜
// 遗传：父母视力
// 用眼环境：时长
// BMI
// 眼压
// 角膜曲率: K1
// 视力：视力表
// 眼轴：长度

// 远视储备：屈光：球镜
// XXX：这个其实不对，需要参考视力和屈光才可以
//
func Range_YuanShiChuBei(age float32) (min, max float32) {
	if age >= 0 && age < 2 {

	} else if age >= 2 && age < 4 {

	} else if age >= 4 && age < 6 {

	} else if age >= 6 && age < 8 {

	} else if age >= 8 && age < 9 {

	} else if age >= 9 && age < 10 {

	} else if age >= 10 && age < 11 {

	} else if age >= 11 && age < 12 {

	} else {

		return 0, 0
	}

	return 0, 0
}

// 遗传
func Range_YiChuan(age float32) (min, max float32) {
	if age >= 0 && age < 2 {

	} else if age >= 2 && age < 4 {

	} else if age >= 4 && age < 6 {

	} else if age >= 6 && age < 8 {

	} else if age >= 8 && age < 9 {

	} else if age >= 9 && age < 10 {

	} else if age >= 10 && age < 11 {

	} else if age >= 11 && age < 12 {

	} else {

		return 0, 0
	}
	return 0, 0
}

// 用眼环境
func Range_YongYanHuanJing(age float32) (min, max float32) {
	if age >= 0 && age < 2 {

	} else if age >= 2 && age < 4 {

	} else if age >= 4 && age < 6 {

	} else if age >= 6 && age < 8 {

	} else if age >= 8 && age < 9 {

	} else if age >= 9 && age < 10 {

	} else if age >= 10 && age < 11 {

	} else if age >= 11 && age < 12 {

	} else {

		return 0, 0
	}

	return 0, 0
}

// 体重BMI
func Range_BMI(age float32) (min, max float32) {
	if age >= 0 && age < 2 {

	} else if age >= 2 && age < 4 {

	} else if age >= 4 && age < 6 {

	} else if age >= 6 && age < 8 {

	} else if age >= 8 && age < 9 {

	} else if age >= 9 && age < 10 {

	} else if age >= 10 && age < 11 {

	} else if age >= 11 && age < 12 {

	} else {

		return 0, 0
	}

	return 0, 0
}

// 眼压
func Range_YanYa(age float32) (min, max float32) {
	if age >= 0 && age < 2 {

	} else if age >= 2 && age < 4 {

	} else if age >= 4 && age < 6 {

	} else if age >= 6 && age < 8 {

	} else if age >= 8 && age < 9 {

	} else if age >= 9 && age < 10 {

	} else if age >= 10 && age < 11 {

	} else if age >= 11 && age < 12 {

	} else {

		return 0, 0
	}

	return 0, 0
}

// 角膜曲率
func Range_JiaoMoQuLv(age float32) (min, max float32) {
	if age >= 0 && age < 2 {

	} else if age >= 2 && age < 4 {

	} else if age >= 4 && age < 6 {

	} else if age >= 6 && age < 8 {

	} else if age >= 8 && age < 9 {

	} else if age >= 9 && age < 10 {

	} else if age >= 10 && age < 11 {

	} else if age >= 11 && age < 12 {

	} else {

		return 0, 0
	}

	return 0, 0
}

// 视力
func Range_ShiLi(age float32) (min, max float32) {
	if age >= 0 && age < 2 {

	} else if age >= 2 && age < 4 {

	} else if age >= 4 && age < 6 {

	} else if age >= 6 && age < 8 {

	} else if age >= 8 && age < 9 {

	} else if age >= 9 && age < 10 {

	} else if age >= 10 && age < 11 {

	} else if age >= 11 && age < 12 {

	} else {

		return 0, 0
	}

	return 0, 0
}

// 眼轴
func Range_YanZhouFaYu(age float32) (min, max float32) {
	if age >= 0 && age < 2 {

	} else if age >= 2 && age < 4 {

	} else if age >= 4 && age < 6 {

	} else if age >= 6 && age < 8 {

	} else if age >= 8 && age < 9 {

	} else if age >= 9 && age < 10 {

	} else if age >= 10 && age < 11 {

	} else if age >= 11 && age < 12 {

	} else {

		return 0, 0
	}

	return 0, 0
}
