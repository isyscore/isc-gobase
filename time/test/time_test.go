package test

import (
	"fmt"
	"testing"
	t0 "time"

	"github.com/isyscore/isc-gobase/time"
)

func TestTime(t *testing.T) {
	t1 := time.TimeInMillis()
	t2 := time.TimeInSeconds()
	t3 := time.TimeInNano()
	t4 := time.TimeInMicro()
	t.Logf("t1: %v, t2: %v", t1, t2)
	t.Logf("t3: %v, t4: %v", t3, t4)
}

func TestTimeBetween(t *testing.T) {
	athen := t0.Date(2000, t0.January, 1, 12, 0, 0, 0, t0.UTC)
	anow := t0.Date(2022, t0.March, 1, 16, 21, 0, 0, t0.UTC)
	t.Logf("years between: %v", time.YearsBetween(anow, athen))
	t.Logf("months between: %v", time.MonthsBetween(anow, athen))
	t.Logf("days between: %v", time.DaysBetween(anow, athen))
	t.Logf("hours between: %v", time.HoursBetween(anow, athen))
	t.Logf("minutes between: %v", time.MinutesBetween(anow, athen))
	t.Logf("seconds between: %v", time.SecondsBetween(anow, athen))
	t.Logf("milliseconds between: %v", time.MilliSecondsBetween(anow, athen))
}

func TestTimeSpan(t *testing.T) {
	athen := t0.Date(2000, t0.January, 1, 12, 0, 0, 0, t0.UTC)
	anow := t0.Date(2022, t0.March, 1, 16, 21, 0, 0, t0.UTC)
	t.Logf("years between: %v", time.YearSpan(anow, athen))
	t.Logf("months between: %v", time.MonthSpan(anow, athen))
	t.Logf("days between: %v", time.DaySpan(anow, athen))
	t.Logf("hours between: %v", time.HourSpan(anow, athen))
	t.Logf("minutes between: %v", time.MinuteSpan(anow, athen))
	t.Logf("seconds between: %v", time.SecondSpan(anow, athen))
	t.Logf("milliseconds between: %v", time.MilliSecondSpan(anow, athen))
}

func TestNumToTimeDuration(t *testing.T) {
	data := time.NumToTimeDuration(3, t0.Hour)
	fmt.Println(data.Milliseconds())
}
