package slick

import (
	"testing"
	"time"

	"github.com/nlopes/slack"
	"github.com/stretchr/testify/assert"
)

func TestShouldReturnTheNextWeekdayTime(t *testing.T) {
	layout := "2006-01-02 15:04:05 -0700 MST"
	now, _ := time.Parse(layout, "2018-10-23 12:00:00 +0000 UTC")

	expectedTs, _ := time.Parse(layout, "2018-10-30 12:00:00 +0000 UTC")

	wday, _ := NextWeekdayTime(now, time.Tuesday, 12, 0)

	assert.NotZero(t, wday)
	assert.Equal(t, wday, expectedTs)
}

func TestShouldReturnTheAskedWeekday(t *testing.T) {
	layout := "2006-01-02 15:04:05 -0700 MST"
	now, _ := time.Parse(layout, "2018-10-23 12:00:00 +0000 UTC")

	wday, _ := NextWeekdayTime(now, time.Sunday, 12, 0)

	assert.Equal(t, wday.Weekday(), time.Sunday)
}

func TestShouldReturnTheAskedHourAndMinute(t *testing.T) {
	layout := "2006-01-02 15:04:05 -0700 MST"
	now, _ := time.Parse(layout, "2018-10-23 12:00:00 +0000 UTC")

	wday, _ := NextWeekdayTime(now, time.Sunday, 12, 0)

	assert.Equal(t, wday.Hour(), 12)
	assert.Equal(t, wday.Minute(), 0)
}

func TestShouldReturnTheDurationInNanoseconds(t *testing.T) {
	layout := "2006-01-02 15:04:05 -0700 MST"
	now, _ := time.Parse(layout, "2018-10-23 12:00:00 +0000 UTC")
	expectedTs, _ := time.Parse(layout, "2018-10-30 12:00:00 +0000 UTC")

	expectedDuration := expectedTs.Sub(now)
	_, duration := NextWeekdayTime(now, time.Tuesday, 12, 0)

	assert.Equal(t, duration, expectedDuration)
}

func TestShouldReturnTheUnixFromTimestamp(t *testing.T) {
	slackTs := "1355517523.000005"
	expectedTs := slack.JSONTime(1355517523)
	unixTs := unixFromTimestamp(slackTs)

	assert.Equal(t, expectedTs, unixTs)
}

func TestShouldReturnCeroIfGivenTimestampIsInvalid(t *testing.T) {
	slackTs := "invalid"
	expectedTs := slack.JSONTime(0)
	unixTs := unixFromTimestamp(slackTs)

	assert.Equal(t, expectedTs, unixTs)
}
