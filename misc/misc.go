package misc

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kamasamikon/miego/klog"
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

// UintTime is to convert time to 20060102150305
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

// Epos convert Excel Position to index, e.g. A .... AZ .... BC to 0, ... 16, ...
func Epos(s string) int {
	old := 0
	for _, c := range s {
		n := c - 'a' + 1
		old = old*26 + int(n)
	}
	return old - 1
}

const MIMEJSON = "application/json;charset=utf-8"

// HTTPPost post json data to peer and convert the response to pongObj structure
func HTTPPost(url string, pingObj interface{}, pongObj interface{}) error {
	var pingString string

	if pingObj == nil {
		pingString = ""
	} else {
		if s, ok := pingObj.(string); ok {
			pingString = s
		} else {
			bytes, ea := json.Marshal(pingObj)
			if ea != nil {
				return ea
			}
			pingString = string(bytes)
		}
	}
	klog.D(pingString)

	r, eb := http.Post(url, MIMEJSON, strings.NewReader(pingString))
	if eb != nil {
		klog.E(eb.Error())
		return eb
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		klog.E("%d", r.StatusCode)
		return fmt.Errorf("StatusCode == %d", r.StatusCode)
	}

	if pongObj == nil {
		return nil
	}
	return json.NewDecoder(r.Body).Decode(pongObj)
}

// HTTPGet convert the response to pongObj structure
func HTTPGet(url string, pongObj interface{}) error {
	r, eb := http.Get(url)
	if eb != nil {
		klog.E(eb.Error())
		return eb
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		klog.E("%d", r.StatusCode)
		return errors.New(fmt.Sprintf("StatusCode == %d", r.StatusCode))
	}

	if pongObj == nil {
		return nil
	}
	return json.NewDecoder(r.Body).Decode(pongObj)
}
