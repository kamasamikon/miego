package misc

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"runtime"
	"sync"
	"time"

	"crypto/md5"
	"encoding/hex"
)

func Shuffle(slice []interface{}) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

func Shuffle2(slice []interface{}) []interface{} {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}

	newSlice := make([]interface{}, len(slice))
	copy(newSlice, slice)

	return newSlice
}

func MD5file(path string) string {
	ctx := md5.New()

	if data, err := ioutil.ReadFile(path); err != nil {
		return ""
	} else {
		ctx.Write(data)
	}

	return hex.EncodeToString(ctx.Sum(nil))
}

func MD5Reader(src io.Reader) string {
	r := bufio.NewReader(src)
	h := md5.New()

	_, err := io.Copy(h, r)
	if err == nil {
		return hex.EncodeToString(h.Sum(nil))
	}
	return ""
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
