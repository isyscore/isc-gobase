package time

import "time"

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
	FmtY        = "2006"
	FmtYYYYMMdd  = "20060102"

	FmtHmsSSSMore = "15:04:05.000000000"
	FmtHmsSSS     = "15:04:05.000"
	FmtHms        = "15:04:05"
	FmtHm         = "15:04"
	FmtH          = "15"
)

func TimeToStringYmdHms(t time.Time) string {
	return t.Format(FmtYMdHms)
}

func TimeToStringYmdHmsS(t time.Time) string {
	return t.Format(FmtYMdHmsSSS)
}

func TimeToStringFormat(t time.Time, format string) string {
	return t.Format(format)
}

func ParseTimeYmsHms(timeStr string) (time.Time, error) {
	return time.ParseInLocation(Year+"-"+Month+"-"+Day+" "+Hour+":"+Minute+":"+Second, timeStr, time.Local)
}

func ParseTimeYmsHmsS(timeStr string) (time.Time, error) {
	return time.ParseInLocation(Year+"-"+Month+"-"+Day+" "+Hour+":"+Minute+":"+Second+".000", timeStr, time.Local)
}

func ParseTimeYmsHmsLoc(timeStr string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(Year+"-"+Month+"-"+Day+" "+Hour+":"+Minute+":"+Second, timeStr, loc)
}

func ParseTimeYmsHmsSLoc(timeStr string, loc *time.Location) (time.Time, error) {
	return time.ParseInLocation(Year+"-"+Month+"-"+Day+" "+Hour+":"+Minute+":"+Second+".000", timeStr, loc)
}
