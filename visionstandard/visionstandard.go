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
