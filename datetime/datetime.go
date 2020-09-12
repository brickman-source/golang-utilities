package datetime

import (
	"github.com/araddon/dateparse"
	"time"
)

type TimeZone string

const (
	TimeZoneShanghai TimeZone = "Asia/Shanghai"
)

func ParseWithTimeZone(layout string, value string, timezone TimeZone) (time.Time, error) {
	l, err := time.LoadLocation(string(timezone))
	if err != nil {
		l = time.Local
	}
	if layout != "" {
		lt, err := time.ParseInLocation(layout, value, l)
		if err == nil {
			return lt, nil
		}
	}
	ret, err := dateparse.ParseIn(value, l)
	return ret, err
}
