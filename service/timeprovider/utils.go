package timeprovider

import "time"

func UTCStringToDate(utcString string) (time.Time, error) {
	return time.Parse(UTCFormat, utcString)
}

func DateToUTCString(date time.Time) string {
	return date.UTC().Format(UTCFormat)
}
