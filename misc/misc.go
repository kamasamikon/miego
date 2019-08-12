package misc

import (
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"
)

const (
	KC_RAND_KIND_NUM = iota
	KC_RAND_KIND_LOWER
	KC_RAND_KIND_UPPER
	KC_RAND_KIND_ALL
)

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

// UintTime : convert time to 20060102150305
func UintTime(t time.Time) uint64 {
	nnnn := t.Year()
	yy := t.Month()
	rr := t.Day()
	ss := t.Hour()
	ff := t.Minute()
	mm := t.Second()

	s := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", nnnn, yy, rr, ss, ff, mm)
	res, _ := strconv.ParseUint(s, 0, 64)

	return res
}

// UintTimeNow : convert time to 20060102150305
func UintTimeNow() uint64 {
	t := time.Now()

	nnnn := t.Year()
	yy := t.Month()
	rr := t.Day()
	ss := t.Hour()
	ff := t.Minute()
	mm := t.Second()

	s := fmt.Sprintf("%04d%02d%02d%02d%02d%02d", nnnn, yy, rr, ss, ff, mm)
	res, _ := strconv.ParseUint(s, 0, 64)

	return res
}

// UintDate : convert time to 20060102
func UintDate(t time.Time) uint {
	nnnn := t.Year()
	yy := t.Month()
	rr := t.Day()

	s := fmt.Sprintf("%04d%02d%02d", nnnn, yy, rr)
	res, _ := strconv.ParseUint(s, 0, 64)

	return uint(res)
}

// UintDateNow : convert time to 20060102
func UintDateNow() uint {
	t := time.Now()

	nnnn := t.Year()
	yy := t.Month()
	rr := t.Day()

	s := fmt.Sprintf("%04d%02d%02d", nnnn, yy, rr)
	res, _ := strconv.ParseUint(s, 0, 64)

	return uint(res)
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

// Atoi : atoi, if fail return default value
func Atoi(a string, def int64) int64 {
	x, e := strconv.ParseInt(a, 0, 64)
	if e != nil {
		return def
	}
	return x
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
