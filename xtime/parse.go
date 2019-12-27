package xtime

import (
	"time"
)

var locLocal *time.Location

func init() {
	locLocal, _ = time.LoadLocation("Local")
}

func Parse(layout string, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, locLocal)
}
