package xtime

import (
	"fmt"
	"strconv"
	"time"
)

// TimeToNum : convert time to 20060102150305
func TimeToNum(t time.Time) uint64 {
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

// TimeNowToNum : convert now time to 20060102150305
func TimeNowToNum() uint64 {
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

// DashTimeToNum : convert now time to 20060102150305
func DashTimeToNum(stime string) uint64 {
	t, err := time.Parse("2006-01-02 15:04:05", stime)
	if err != nil {
		return 0
	}

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

// DateToNum : convert date to 20060102
func DateToNum(t time.Time) uint {
	nnnn := t.Year()
	yy := t.Month()
	rr := t.Day()

	s := fmt.Sprintf("%04d%02d%02d", nnnn, yy, rr)
	res, _ := strconv.ParseUint(s, 0, 64)

	return uint(res)
}

// DateNowToNum : convert now date to 20060102
func DateNowToNum() uint {
	t := time.Now()

	nnnn := t.Year()
	yy := t.Month()
	rr := t.Day()

	s := fmt.Sprintf("%04d%02d%02d", nnnn, yy, rr)
	res, _ := strconv.ParseUint(s, 0, 64)

	return uint(res)
}

// DashDateToNum : convert time to 20060102
func DashDateToNum(date string) uint {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return 0
	}

	nnnn := t.Year()
	yy := t.Month()
	rr := t.Day()

	s := fmt.Sprintf("%04d%02d%02d", nnnn, yy, rr)
	res, _ := strconv.ParseUint(s, 0, 64)

	return uint(res)
}

func NumTimeToTime(num uint64) (time.Time, error) {
	return Parse("20060102150405", fmt.Sprintf("%d", num))
}

func NumDateToTime(num uint64) (time.Time, error) {
	return Parse("20060102", fmt.Sprintf("%d", num))
}

// NumDateToDash : 20060102 -> 2006-01-02
func NumDateToDash(date uint) string {
	// XXX: Parser will check bad time
	t, err := time.Parse("20060102", fmt.Sprintf("%d", date))
	if err != nil {
		return ""
	}

	nnnn := t.Year()
	yy := t.Month()
	rr := t.Day()

	return fmt.Sprintf("%04d-%02d-%02d", nnnn, yy, rr)
}

// StrDateToDash : 20060102 -> 2006-01-02
func StrDateToDash(date string) string {
	// XXX: Parser will check bad time
	t, err := time.Parse("20060102", date)
	if err != nil {
		return ""
	}

	nnnn := t.Year()
	yy := t.Month()
	rr := t.Day()

	return fmt.Sprintf("%04d-%02d-%02d", nnnn, yy, rr)
}
