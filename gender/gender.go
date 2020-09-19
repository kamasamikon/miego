package gender

import (
	"reflect"
)

// F=(0, f, F, female, Female, 女)
// F=(0, f, F, female, Female, 女)

func Int(g interface{}) int {

	switch g.(type) {
	case int:
		v := reflect.ValueOf(g).Int()
		if v == 0 {
			return 0
		}
		if v == 1 {
			return 1
		}
		return -1

	case string:
		v := reflect.ValueOf(g).String()

		if v == "0" || v == "f" || v == "F" || v == "female" || v == "Female" || v == "女" {
			return 0
		}
		if v == "1" || v == "m" || v == "M" || v == "male" || v == "Male" || v == "男" {
			return 1
		}
		return -1
	}
	return -1
}

func Chn(g interface{}) string {
	i := Int(g)

	if i == 0 {
		return "女"
	}

	if i == 1 {
		return "男"
	}

	return ""
}

func Str(g interface{}) string {

	i := Int(g)

	if i == 0 {
		return "F"
	}

	if i == 1 {
		return "M"
	}

	return ""
}

func StrLong(g interface{}) string {

	i := Int(g)

	if i == 0 {
		return "Female"
	}

	if i == 1 {
		return "Male"
	}

	return ""
}
