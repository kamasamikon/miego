package misc

import (
	"fmt"
	"strconv"
	"time"
)

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

// UintTimeNow : convert now time to 20060102150305
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

// UintDate : convert date to 20060102
func UintDate(t time.Time) uint {
	nnnn := t.Year()
	yy := t.Month()
	rr := t.Day()

	s := fmt.Sprintf("%04d%02d%02d", nnnn, yy, rr)
	res, _ := strconv.ParseUint(s, 0, 64)

	return uint(res)
}

// UintDateNow : convert now date to 20060102
func UintDateNow() uint {
	t := time.Now()

	nnnn := t.Year()
	yy := t.Month()
	rr := t.Day()

	s := fmt.Sprintf("%04d%02d%02d", nnnn, yy, rr)
	res, _ := strconv.ParseUint(s, 0, 64)

	return uint(res)
}

// UintDateStr : convert time to 20060102
func UintDateStr(date string) uint {
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
