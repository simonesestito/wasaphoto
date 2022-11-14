package timeprovider

import "time"

type RealTimeProvider struct{}

func (RealTimeProvider) Now() time.Time {
	return time.Now()
}

func (time RealTimeProvider) UTCString() string {
	return time.Now().UTC().Format(UTCFormat)
}
