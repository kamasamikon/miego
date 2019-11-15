package in

// S : Check "aaa" in ["aaa", "bbb", "..."]
func S(s string, arr ...string) bool {
	for _, i := range arr {
		if s == i {
			return true
		}
	}
	return false
}
