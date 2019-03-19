package util

import (
	"time"
)

func GetZeroBetween(begin int64, end int64) uint32 {
	if end <= begin {
		return 0
	}

	return (uint32)((end+8*60*60)/(24*60*60) - (begin+8*60*60)/(24*60*60))
}

func GetWeekdayBetween(begin int64, end int64) uint32 {
	if end <= begin {
		return 0
	}

	b := GetWeekdayStart(begin)
	e := GetWeekdayStart(end)

	w := uint32((e - b) / (7 * 24 * 3600))

	return w
}

func GetDayStart(t int64) int64 {
	if t < 86400 {
		return 0
	}
	_0 := t/86400*86400 - 28800
	if _0 <= (t - 86400) {
		_0 += 86400
	}
	return _0
}

func GetWeekdayStart(t int64) int64 {

	t = GetDayStart(t)
	tm := time.Unix(t, 0)

	if uint32(tm.Weekday()) != 0 {
		t -= int64((tm.Weekday() - 1) * 24 * 60 * 60)
	} else {
		t -= 6 * 24 * 60 * 60
	}

	return t
}
