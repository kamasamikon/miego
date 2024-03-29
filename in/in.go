package in

// S : string, Check "aaa" in ["aaa", "bbb", "..."]
func S(s string, arr ...string) bool {
	for _, i := range arr {
		if s == i {
			return true
		}
	}
	return false
}

// I : Number, Check 2 in [3, 4, 5, ...]
func I(i int, arr ...int) bool {
	for _, t := range arr {
		if i == t {
			return true
		}
	}
	return false
}

// C : Char, Check 2 in [3, 4, 5, ...]
func C(i byte, arr ...byte) bool {
	for _, t := range arr {
		if i == t {
			return true
		}
	}
	return false
}
