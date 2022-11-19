package utils

import (
	"github.com/simonesestito/wasaphoto/service/timeprovider"
	"time"
)

func UTCStringToDate(utcString string) (time.Time, error) {
	return time.Parse(timeprovider.UTCFormat, utcString)
}
