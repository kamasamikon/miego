package misc

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

const (
	KC_RAND_KIND_NUM = iota
	KC_RAND_KIND_LOWER
	KC_RAND_KIND_UPPER
	KC_RAND_KIND_ALL
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

// Xrand : kind: "nlupa"
func Xrand(size int, kind string) []byte {
	num := 0

	//
	// Parse kind
	//
	for c := range kind {
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

	if num == 0 {
		num = NUM
	}

	//
	// Load Cache
	//
	str, ok := charMap[num]
	if !ok {
		tmp := ""

		if num%NUM != 0 {
			tmp += arrNUM
		}

		if num%LOW != 0 {
			tmp += arrLOW
		}

		if num%UPP != 0 {
			tmp += arrUPP
		}

		if num%PUN != 0 {
			tmp += arrPUN
		}

		charMap[num] = tmp
		str = tmp
	}

	result := make([]byte, size)
	arrSize := len(str)
	for i := 0; i < size; i++ {
		index := rand.Intn(arrSize)
		result[i] = uint8(str[index])
	}

	return result
}

func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

func ReverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

func ReverseBytes(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

func MemConsume() {
	var c chan int
	var wg sync.WaitGroup
	const goroutineNum = 1e4

	memConsumed := func() uint64 {
		runtime.GC() //GC，排除对象影响

		var memStat runtime.MemStats
		runtime.ReadMemStats(&memStat)
		return memStat.Sys
	}
	noop := func() {
		wg.Done()
		<-c //防止goroutine退出，内存被释放
	}

	wg.Add(goroutineNum)
	before := memConsumed() //获取创建goroutine前内存
	for i := 0; i < goroutineNum; i++ {
		go noop()
	}
	wg.Wait()

	after := memConsumed() //获取创建goroutine后内存
	fmt.Printf("%.3f KB\n", float64(after-before)/goroutineNum/1000)
}

// Epos : convert Excel Position to index, e.g. A .... AZ .... BC to 0, ... 16, ...
func Epos(s string) int {
	old := 0
	for _, c := range s {
		n := c - 'a' + 1
		old = old*26 + int(n)
	}
	return old - 1
}
