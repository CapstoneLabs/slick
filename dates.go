package slick

import (
	"strconv"
	"strings"
	"time"

	"github.com/nlopes/slack"
)

// NextWeekdayTime from a given Time and a Weekday + hour + minute
// returns the next WeekdayTime and the elapsed time between the two instants
// as an int64 nanosecond count.
func NextWeekdayTime(from time.Time, w time.Weekday, hour, min int) (time.Time, time.Duration) {
	t := from.UTC()
	nowWeekday := t.Weekday()
	nowYear, nowMonth, nowDay := t.Date()

	delta := (int(w) - int(nowWeekday) + 7) % 7
	res := time.Date(nowYear, nowMonth, nowDay+delta, hour, min, 0, 0, t.Location())

	if res.Sub(t) <= 0 {
		res = res.AddDate(0, 0, 7)
	}
	return res, res.Sub(t)
}

// AfterNextWeekdayTime from a given Time and a Weekday + hour + minute
// returns a channel that waits for the duration to elapse and then sends
// the current time on the returned channel.
func AfterNextWeekdayTime(from time.Time, w time.Weekday, hour, min int) <-chan time.Time {
	_, duration := NextWeekdayTime(from, w, hour, min)
	return time.After(duration)
}

// unixFromTimestamp from an slack unique (per-channel) timestamp
// returns the unix timestamp.
func unixFromTimestamp(ts string) slack.JSONTime {
	i, _ := strconv.Atoi(strings.Split(ts, ".")[0])
	return slack.JSONTime(i)
}
