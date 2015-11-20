package parse

import (
	"strings"
	"time"
)

func GetMonth() string {
	t := time.Now().Unix()
	return ParseTime(t, "200601")
}

func GetNextMonth() string {
	t := time.Now().Unix() + 24*3600*2
	return ParseTime(t, "200601")
}

func IsLastDayOfMonth() bool {
	t := time.Now().Unix()
	month := ParseTime(t, "01")

	nt := t + 24*3600*2
	nmonth := ParseTime(nt, "01")
	if month != nmonth {
		return true
	}
	return false
}

func DateStrtoUnixStr(date string) string {
	if date == "" {
		date = ParseTime(time.Now().Unix(), "2006-01-02 15:04:05")
	}
	date = strings.Replace(date, " ", "T", 1)
	return date + "+08:00"
}

func ParseDate(date, format string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, _ := time.ParseInLocation(format, date, loc)
	return t.Unix()
}

func ParseTime(t int64, format string) string {
	return time.Unix(t, 0).Format(format)
}

func ParseWeekday(date string, format string) string {
	t, _ := time.Parse(format, date)
	return t.Weekday().String()
}

func DateFormatTrans(date, preformat, format string) string {
	t, _ := time.Parse(preformat, date)
	return t.Format(format)
}

func UTCDateToLocalDate(date string) string {
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", date, time.UTC)
	return time.Unix(t.Unix(), 0).Format("2006-01-02 15:04:05")
}
