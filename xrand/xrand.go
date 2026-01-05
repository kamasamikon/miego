package xrand

import (
	"math/rand"
)

const (
	NUM = 2
	LOW = 3
	UPP = 5
	PUN = 7
	ALL = 11
)

var arrNUM = "0123456789"
var arrLOW = "abcdefghijklmnopqrstuvwxyz"
var arrUPP = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
var arrPUN = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

var charMap map[int]string = make(map[int]string)

func charMapGet(num int) (string, int) {
	//
	// Load Cache
	//
	str, ok := charMap[num]
	if !ok {
		tmp := ""

		if num%NUM == 0 {
			tmp += arrNUM
		}

		if num%LOW == 0 {
			tmp += arrLOW
		}

		if num%UPP == 0 {
			tmp += arrUPP
		}

		if num%PUN == 0 {
			tmp += arrPUN
		}

		charMap[num] = tmp
		str = tmp
	}

	return str, len(str)
}

// Range : kind: "nlupa"
func Rand(size int, kind string) []byte {
	num := 1

	//
	// Parse kind
	//
	for _, c := range kind {
		if c == 'n' {
			if num%NUM != 0 {
				num *= NUM
			}
		} else if c == 'l' {
			if num%LOW != 0 {
				num *= LOW
			}
		} else if c == 'u' {
			if num%UPP != 0 {
				num *= UPP
			}
		} else if c == 'p' {
			if num%PUN != 0 {
				num *= PUN
			}
		} else if c == 'a' {
			num = NUM * LOW * UPP * PUN
		}
	}

	if num == 1 {
		num = NUM
	}

	result := make([]byte, size)

	//
	// Load Cache
	//
	str, strlen := charMapGet(num)
	for i := 0; i < size; i++ {
		index := rand.Intn(strlen)
		result[i] = uint8(str[index])
	}

	return result
}

func Num(size int) []byte {
	result := make([]byte, size)

	str, strlen := charMapGet(NUM)
	for i := 0; i < size; i++ {
		index := rand.Intn(strlen)
		result[i] = uint8(str[index])
	}

	return result
}
