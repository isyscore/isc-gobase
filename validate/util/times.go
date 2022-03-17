package util

import (
	"fmt"
	"github.com/isyscore/isc-gobase/logger"
	"regexp"
	"strings"
	"time"
)

var yRegex = regexp.MustCompile("^(\\d){4}$")
var ymRegex = regexp.MustCompile("^(\\d){4}-(\\d){2}$")
var ymdRegex = regexp.MustCompile("^(\\d){4}-(\\d){2}-(\\d){2}$")
var ymdHRegex = regexp.MustCompile("^(\\d){4}-(\\d){2}-(\\d){2} (\\d){2}$")
var ymdHmRegex = regexp.MustCompile("^(\\d){4}-(\\d){2}-(\\d){2} (\\d){2}:(\\d){2}$")
var ymdHmsRegex = regexp.MustCompile("^(\\d){4}-(\\d){2}-(\\d){2} (\\d){2}:(\\d){2}:(\\d){2}$")
var ymdHmsSRegex = regexp.MustCompile("^(\\d){4}-(\\d){2}-(\\d){2} (\\d){2}:(\\d){2}:(\\d){2}.(\\d){3}$")

var yTime = "2006"
var ymTime = "2006-01"
var ymdTime = "2006-01-02"
var ymdHTime = "2006-01-02 15"
var ymdHmTime = "2006-01-02 15:04"
var ymdHmsTime = "2006-01-02 15:04:05"
var ymdHmsSTime = "2006-01-02 15:04:05.000"

var EmptyTime = time.Time{}

func AddHour(times time.Time, plusOrMinus string, seconds string) time.Time {
	h, _ := time.ParseDuration(fmt.Sprintf("%s%v", plusOrMinus, seconds))
	return times.Add(h)
}

func AddMinutes(times time.Time, plusOrMinus string, minutes string) time.Time {
	h, _ := time.ParseDuration(fmt.Sprintf("%s%v", plusOrMinus, minutes))
	return times.Add(h)
}

func AddSeconds(times time.Time, plusOrMinus string, hours string) time.Time {
	h, _ := time.ParseDuration(fmt.Sprintf("%s%v", plusOrMinus, hours))
	return times.Add(h)
}

func AddDays(times time.Time, days int) time.Time {
	return times.AddDate(0, 0, days)
}

func AddMonths(times time.Time, month int) time.Time {
	return times.AddDate(0, month, 0)
}

func AddYears(times time.Time, year int) time.Time {
	return times.AddDate(year, 0, 0)
}

func ParseTime(timeStr string) time.Time {
	timeStr = strings.TrimSpace(timeStr)
	timeStr = strings.TrimSpace(strings.ReplaceAll(timeStr, "\\'", " "))

	if timeStr == "" {
		return EmptyTime
	}
	if yRegex.MatchString(timeStr) {
		if times, err := time.Parse(yTime, timeStr); err == nil {
			return times
		} else {
			logger.Error("解析时间错误, err: %v", err)
		}
	} else if ymRegex.MatchString(timeStr) {
		if times, err := time.Parse(ymTime, timeStr); err == nil {
			return times
		} else {
			logger.Error("解析时间错误, err: %v", err)
		}
	} else if ymdRegex.MatchString(timeStr) {
		if times, err := time.Parse(ymdTime, timeStr); err == nil {
			return times
		} else {
			logger.Error("解析时间错误, err: %v", err)
		}
	} else if ymdHRegex.MatchString(timeStr) {
		if times, err := time.Parse(ymdHTime, timeStr); err == nil {
			return times
		} else {
			logger.Error("解析时间错误, err: %v", err)
		}
	} else if ymdHmRegex.MatchString(timeStr) {
		if times, err := time.Parse(ymdHmTime, timeStr); err == nil {
			return times
		} else {
			logger.Error("解析时间错误, err: %v", err)
		}
	} else if ymdHmsRegex.MatchString(timeStr) {
		if times, err := time.Parse(ymdHmsTime, timeStr); err == nil {
			return times
		} else {
			logger.Error("解析时间错误, err: %v", err)
		}
	} else if ymdHmsSRegex.MatchString(timeStr) {
		if times, err := time.Parse(ymdHmsSTime, timeStr); err == nil {
			return times
		} else {
			logger.Error("解析时间错误, err: %v", err)
		}
	} else {
		logger.Error("解析时间错误, time: %v", timeStr)
	}
	return EmptyTime
}

func IsTimeEmpty(time time.Time) bool {
	if time == EmptyTime {
		return true
	}
	return false
}
