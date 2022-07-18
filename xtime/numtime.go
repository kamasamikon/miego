package xtime

import (
	"fmt"
	"strconv"
	"time"
)

// TimeToNum : convert time to 20060102150305
func TimeToNum(t time.Time) uint64 {
	res, _ := strconv.ParseUint(t.Format("20060102150405"), 0, 64)
	return res
}

// TimeNowToNum : convert now time to 20060102150305
func TimeNowToNum() uint64 {
	res, _ := strconv.ParseUint(time.Now().Format("20060102150405"), 0, 64)
	return res
}

// DashTimeToNum : convert now time to 20060102150305
func DashTimeToNum(stime string) uint64 {
	t, err := time.Parse("2006-01-02 15:04:05", stime)
	if err != nil {
		return 0
	}
	res, _ := strconv.ParseUint(t.Format("20060102150405"), 0, 64)
	return res
}

// DateToNum : convert date to 20060102
func DateToNum(t time.Time) int {
	res, _ := strconv.ParseInt(t.Format("20060102"), 0, 64)
	return int(res)
}

// DateNowToNum : convert now date to 20060102
func DateNowToNum() int {
	res, _ := strconv.ParseInt(time.Now().Format("20060102"), 0, 64)
	return int(res)
}

// DashDateToNum : convert time to 20060102
func DashDateToNum(date string) int {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0
	}
	res, _ := strconv.ParseInt(t.Format("20060102"), 0, 64)
	return int(res)
}

func NumTimeToTime(num uint64) (time.Time, error) {
	return Parse("20060102150405", fmt.Sprintf("%d", num))
}

func NumDateToTime(num uint64) (time.Time, error) {
	return Parse("20060102", fmt.Sprintf("%d", num))
}

// NumDateToDash : 20060102 -> 2006-01-02
func NumDateToDash(date int) string {
	return StrDateToDash(fmt.Sprintf("%d", date))
}

// StrDateToDash : 20060102 -> 2006-01-02
func StrDateToDash(date string) string {
	t, err := time.Parse("20060102", date)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%04d-%02d-%02d", t.Year(), t.Month(), t.Day())
}
