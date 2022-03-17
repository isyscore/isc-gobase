package time

import (
	"math"
	t0 "time"
)

var (
	Year   = "2006"
	Month  = "01"
	Day    = "02"
	Hour   = "15"
	Minute = "04"
	Second = "05"

	FmtYMdHmsSSS = "2006-01-02 15:04:05.000"
	FmtYMdHmsS   = "2006-01-02 15:04:05.0"
	FmtYMdHms    = "2006-01-02 15:04:05"
	FmtYMdHm     = "2006-01-02 15:04"
	FmtYMdH      = "2006-01-02 15"
	FmtYMd       = "2006-01-02"
	FmtYM        = "2006-01"
	FmtY         = "2006"
	FmtYYYYMMdd  = "20060102"

	FmtHmsSSSMore = "15:04:05.000000000"
	FmtHmsSSS     = "15:04:05.000"
	FmtHms        = "15:04:05"
	FmtHm         = "15:04"
	FmtH          = "15"
)

type FPCDateTime float64

const (
	DateDelta   = 693594 // Days between 1/1/0001 and 12/31/1899
	HoursPerDay = 24
	MinsPerHour = 60
	SecsPerMin  = 60
	MSecsPerSec = 1000

	MinsPerDay  = HoursPerDay * MinsPerHour
	SecsPerHour = SecsPerMin * MinsPerHour
	SecsPerDay  = MinsPerDay * SecsPerMin
	MSecsPerDay = SecsPerDay * MSecsPerSec

	OneMillisecond  = FPCDateTime(1) / MSecsPerDay
	HalfMilliSecond = OneMillisecond / 2

	JulianEpoch = FPCDateTime(-2415018.5)
	UnixEpoch   = JulianEpoch + FPCDateTime(2440587.5)

	ApproxDaysPerMonth = 30.4375
	ApproxDaysPerYear  = 365.25
)

func TimeToStringYmdHms(t t0.Time) string {
	return t.Format(FmtYMdHms)
}

func TimeToStringYmdHmsS(t t0.Time) string {
	return t.Format(FmtYMdHmsSSS)
}

func TimeToStringFormat(t t0.Time, format string) string {
	return t.Format(format)
}

func ParseTimeYmsHms(timeStr string) (t0.Time, error) {
	return t0.ParseInLocation(Year+"-"+Month+"-"+Day+" "+Hour+":"+Minute+":"+Second, timeStr, t0.Local)
}

func ParseTimeYmsHmsS(timeStr string) (t0.Time, error) {
	return t0.ParseInLocation(Year+"-"+Month+"-"+Day+" "+Hour+":"+Minute+":"+Second+".000", timeStr, t0.Local)
}

func ParseTimeYmsHmsLoc(timeStr string, loc *t0.Location) (t0.Time, error) {
	return t0.ParseInLocation(Year+"-"+Month+"-"+Day+" "+Hour+":"+Minute+":"+Second, timeStr, loc)
}

func ParseTimeYmsHmsSLoc(timeStr string, loc *t0.Location) (t0.Time, error) {
	return t0.ParseInLocation(Year+"-"+Month+"-"+Day+" "+Hour+":"+Minute+":"+Second+".000", timeStr, loc)
}

// TimeInMillis 13位java时间戳
func TimeInMillis() int64 {
	return t0.Now().UnixMilli()
}

// TimeInSeconds 10位Unix时间戳
func TimeInSeconds() int64 {
	return t0.Now().Unix()
}

// TimeInMicro 16位时间戳
func TimeInMicro() int64 {
	return t0.Now().UnixMicro()
}

// TimeInNano 19位时间戳
func TimeInNano() int64 {
	return t0.Now().UnixNano()
}

//====================================================================================
// 时间日期函数
//====================================================================================

func CurrentMinuteOfDay() int {
	t := t0.Now()
	return t.Hour()*60 + t.Minute()
}

func CurrentSecondOfDay() int {
	t := t0.Now()
	return t.Hour()*3600 + t.Minute()*60 + t.Second()
}

func MinuteOfDay(t t0.Time) int {
	return t.Hour()*60 + t.Minute()
}

func SecondOfDay(t t0.Time) int {
	return t.Hour()*3600 + t.Minute()*60 + t.Second()
}

func MinutesToTime(minutes int) (hour int, minute int) {
	h0 := minutes / 60
	m0 := minutes % 60
	return h0, m0
}

func SecondsToTime(seconds int) (hour int, minute int, second int) {
	h0, m0 := MinutesToTime(seconds / 60)
	s0 := seconds % 60
	return h0, m0, s0
}

func IsLeapYear(year int) bool {
	return (year%4 == 0) && ((year%100 != 0) || (year%400 == 0))
}

func YearsBetween(now, then t0.Time) int {
	return int(math.Trunc(now.Sub(then).Hours() / 24 / ApproxDaysPerYear))
}

func MonthsBetween(now, then t0.Time) int {
	return int(math.Trunc(now.Sub(then).Hours() / 24 / ApproxDaysPerMonth))
}

func DaysBetween(now, then t0.Time) int {
	return int(math.Trunc(now.Sub(then).Hours() / 24))
}

func HoursBetween(now, then t0.Time) int {
	return int(math.Trunc(now.Sub(then).Hours()))
}

func MinutesBetween(now, then t0.Time) int {
	return int(math.Trunc(now.Sub(then).Minutes()))
}

func SecondsBetween(now, then t0.Time) int {
	return int(math.Trunc(now.Sub(then).Seconds()))
}

func MilliSecondsBetween(now, then t0.Time) int64 {
	return now.Sub(then).Milliseconds()
}

func WithInPastYears(now, then t0.Time, years int) bool {
	return YearsBetween(now, then) <= years
}

func WithInPastMonths(now, then t0.Time, months int) bool {
	return MonthsBetween(now, then) <= months
}

func WithInPastDays(now, then t0.Time, days int) bool {
	return DaysBetween(now, then) <= days
}

func WithInPastHours(now, then t0.Time, hours int) bool {
	return HoursBetween(now, then) <= hours
}

func WithInPastMinutes(now, then t0.Time, minutes int) bool {
	return MinutesBetween(now, then) <= minutes
}

func WithInPastSeconds(now, then t0.Time, seconds int) bool {
	return SecondsBetween(now, then) <= seconds
}

func WithInPastMilliSeconds(now, then t0.Time, milliSeconds int64) bool {
	return MilliSecondsBetween(now, then) <= milliSeconds
}

func YearSpan(now, then t0.Time) float64 {
	return now.Sub(then).Hours() / 24 / ApproxDaysPerYear
}

func MonthSpan(now, then t0.Time) float64 {
	return now.Sub(then).Hours() / 24 / ApproxDaysPerMonth
}

func DaySpan(now, then t0.Time) float64 {
	return now.Sub(then).Hours() / 24
}

func HourSpan(now, then t0.Time) float64 {
	return now.Sub(then).Hours()
}

func MinuteSpan(now, then t0.Time) float64 {
	return now.Sub(then).Minutes()
}

func SecondSpan(now, then t0.Time) float64 {
	return now.Sub(then).Seconds()
}

func MilliSecondSpan(now, then t0.Time) int64 {
	return now.Sub(then).Milliseconds()
}

func Now() t0.Time {
	return t0.Now()
}
